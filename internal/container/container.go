package container

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/mihazzz123/m3zold-server/internal/config"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http/handlers"
	"github.com/mihazzz123/m3zold-server/internal/domain/services"
	"github.com/mihazzz123/m3zold-server/internal/infrastructure/repository"
	infrastructure_services "github.com/mihazzz123/m3zold-server/internal/infrastructure/services"
	authusecase "github.com/mihazzz123/m3zold-server/internal/usecase/auth"
	deviceusecase "github.com/mihazzz123/m3zold-server/internal/usecase/device"
	healthusecase "github.com/mihazzz123/m3zold-server/internal/usecase/health"
	userusecase "github.com/mihazzz123/m3zold-server/internal/usecase/user"
)

type Container struct {
	// Core dependencies
	Config *config.Config
	Logger *logrus.Logger
	DB     *pgxpool.Pool

	// Repositories
	UserRepo              *repository.UserRepository
	DeviceRepo            *repository.DeviceRepository
	VerificationEmailRepo *repository.VerificationEmailRepository
	HealthRepo            *repository.HealthRepository
	AuthRepo              *repository.AuthRepository

	// Services
	PasswordService       services.PasswordService
	IDService             services.IDService
	EmailValidatorService services.EmailValidatorService
	UserFactory           services.UserFactory
	TokenService          services.TokenService
	JWTService            services.JWTService

	// Use Cases
	ProfileUseCase *userusecase.ProfileUseCase
	HealthUseCase  *healthusecase.HealthUseCase
	AuthUseCase    *authusecase.AuthUseCase

	// Handlers
	UserHandler              *handlers.UserHandler
	VerificationEmailHandler *handlers.VerificationEmailHandler
	DeviceHandler            *handlers.DeviceHandler
	HealthHandler            *handlers.HealthHandler
	AuthHandler              *handlers.AuthHandler
}

func New(ctx context.Context) (*Container, error) {
	// Initialize config
	cfg := config.New()

	// Initialize logger
	logger := setupLogger(cfg)

	// Initialize database
	db, err := setupDatabase(ctx, cfg.Database.Url, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}

	// Initialize services
	passwordService := infrastructure_services.NewPasswordService(0)
	idService := infrastructure_services.NewIDService()
	emailValidatorService := infrastructure_services.NewEmailValidatorService()
	userFactory := infrastructure_services.NewUserFactory()
	tokenService := infrastructure_services.NewTokenService()
	jwtService := infrastructure_services.NewJWTService(
		cfg.Auth.JWTSecret,
		"m3zold-server",
		"m3zold-client",
	)

	// Initialize repositories
	userRepo := repository.NewUserRepo(db)
	deviceRepo := repository.NewDeviceRepo(db)
	verificationEmailRepo := repository.NewVerificationEmailRepository(db)
	healthRepo := repository.NewHealthRepo(db)
	authRepo := repository.NewAuthRepository(db)

	userService := services.NewUserService(userRepo)

	// UseCases
	profileUseCase := userusecase.NewProfileUseCase(userRepo, emailValidatorService)
	createDeviceUC := deviceusecase.NewCreateUseCase(deviceRepo)
	deleteUseCase := deviceusecase.NewDeleteUseCase(deviceRepo)
	findUseCase := deviceusecase.NewFindUseCase(deviceRepo)
	listDeviceUC := deviceusecase.NewListUseCase(deviceRepo)
	updateStatusUseCase := deviceusecase.NewUpdateStatusUseCase(deviceRepo)
	healthUseCase := healthusecase.NewHealthUseCase(healthRepo)
	authUseCase := authusecase.NewAuthUseCase(
		authRepo,
		userRepo,
		passwordService,
		tokenService,
		jwtService,
		emailValidatorService,
		idService,
		userFactory,
	)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(profileUseCase)
	deviceHandler := handlers.NewDeviceHandler(createDeviceUC, listDeviceUC, findUseCase, updateStatusUseCase, deleteUseCase)
	verificationEmailHandler := handlers.NewVerificationEmailHandler(nil)
	healthHandler := handlers.NewHealthHandler(healthUseCase)
	authHandler := handlers.NewAuthHandler(authUseCase, userService)

	return &Container{
		Config:                cfg,
		Logger:                logger,
		DB:                    db,
		UserRepo:              userRepo,
		DeviceRepo:            deviceRepo,
		VerificationEmailRepo: verificationEmailRepo,
		HealthRepo:            healthRepo,
		AuthRepo:              authRepo,

		PasswordService:       passwordService,
		IDService:             idService,
		EmailValidatorService: emailValidatorService,
		UserFactory:           userFactory,
		TokenService:          tokenService,
		JWTService:            jwtService,

		ProfileUseCase: profileUseCase,
		HealthUseCase:  healthUseCase,
		AuthUseCase:    authUseCase,

		UserHandler:              userHandler,
		DeviceHandler:            deviceHandler,
		VerificationEmailHandler: verificationEmailHandler,
		HealthHandler:            healthHandler,
		AuthHandler:              authHandler,
	}, nil
}

// Close корректно закрывает все ресурсы контейнера
func (c *Container) Close() {
	c.Logger.Info("🔄 Closing container resources...")

	// Закрываем соединение с БД
	if c.DB != nil {
		c.Logger.Info("📦 Closing database connection...")
		c.DB.Close()
		c.Logger.Info("✅ Database connection closed")
	}

	// Закрываем файловые дескрипторы логгера если пишем в файл
	if c.Logger != nil && c.Config.Logger.Output == "file" {
		if file, ok := c.Logger.Out.(*os.File); ok && file != os.Stdout {
			c.Logger.Info("📝 Closing log file...")
			file.Close()
			c.Logger.Info("✅ Log file closed")
		}
	}

	c.Logger.Info("✅ All container resources closed")
}

// GetDB возвращает пул соединений с БД (для миграций и т.д.)
func (c *Container) GetDB() *pgxpool.Pool {
	return c.DB
}

// GetLogger возвращает логгер
func (c *Container) GetLogger() *logrus.Logger {
	return c.Logger
}

// GetConfig возвращает конфиг
func (c *Container) GetConfig() *config.Config {
	return c.Config
}

func setupLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	// Форматтер
	if cfg.Logger.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	}

	// Output
	if cfg.Logger.Output == "stdout" {
		logger.SetOutput(os.Stdout)
	} else if cfg.Logger.Output == "file" {
		file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.SetOutput(os.Stdout)
			logger.Warn("Failed to log to file, using stdout")
		} else {
			logger.SetOutput(file)
		}
	} else {
		logger.SetOutput(os.Stdout)
	}

	// Уровень логирования
	switch cfg.Logger.Level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}

// setupDatabase ожидает подключения к БД с повторными попытками
func setupDatabase(ctx context.Context, dbURL string, logger *logrus.Logger) (*pgxpool.Pool, error) {
	maxAttempts := 10
	retryDelay := 3 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		logger.Infof("Attempt %d/%d to connect to database...", attempt, maxAttempts)

		// Создаем конфиг пула
		config, err := pgxpool.ParseConfig(dbURL)
		if err != nil {
			logger.Warnf("Failed to parse DB config: %v", err)
			if attempt < maxAttempts {
				logger.Infof("Retrying in %v...", retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("failed to parse database config: %w", err)
		}

		// Настраиваем пул соединений
		config.MaxConns = 10
		config.MinConns = 2
		config.HealthCheckPeriod = 1 * time.Minute
		config.MaxConnLifetime = 1 * time.Hour

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			logger.Warnf("Connection attempt %d failed: %v", attempt, err)
			if attempt < maxAttempts {
				logger.Infof("Retrying in %v...", retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("failed to create connection pool: %w", err)
		}

		// Проверяем подключение
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			logger.Warnf("Ping attempt %d failed: %v", attempt, err)
			if attempt < maxAttempts {
				logger.Infof("Retrying in %v...", retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("database ping failed: %w", err)
		}

		logger.Info("✅ Database connection successful")
		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts", maxAttempts)
}
