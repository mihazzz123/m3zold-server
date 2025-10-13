package container

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mihazzz123/m3zold-server/internal/config"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http/handlers"
	"github.com/mihazzz123/m3zold-server/internal/infrastructure/repository"
	infrastructure_services "github.com/mihazzz123/m3zold-server/internal/infrastructure/services"
	healthusecase "github.com/mihazzz123/m3zold-server/internal/usecase"
	"github.com/mihazzz123/m3zold-server/internal/usecase/device"
	userusecase "github.com/mihazzz123/m3zold-server/internal/usecase/user"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Config *config.Config
	Logger *logrus.Logger
	DB     *pgxpool.Pool

	// Health
	HealthUseCase *healthusecase.HealthUseCase
	HealthHandler *handlers.HealthHandler

	UserHandler   *handlers.UserHandler
	DeviceHandler *handlers.DeviceHandler
}

func New(ctx context.Context) (*Container, error) {
	cfg := config.New()
	logger := loggerSetup(cfg)
	// Сначала проверяем и ждем подключение к БД
	db, err := waitForDatabase(ctx, cfg.Database.Url, logger)
	if err != nil {
		logger.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	// Services
	passwordService := infrastructure_services.NewPasswordService(0)
	idService := infrastructure_services.NewIDService()
	emailService := infrastructure_services.NewEmailService()
	userFactory := infrastructure_services.NewUserFactory()

	// Repositories
	userRepo := repository.NewUserRepo(db)
	healthRepo := repository.NewHealthRepo(db)
	deviceRepo := repository.NewDeviceRepo(db)

	// Use Cases
	registerUseCase := userusecase.NewRegisterUseCase(
		userRepo,
		passwordService,
		idService,
		emailService,
		userFactory,
	)

	healthUseCase := healthusecase.NewHealthUseCase(healthRepo)
	healthHandler := handlers.NewHealthHandler(healthUseCase)

	// Device UseCases
	createDeviceUC := device.NewCreateUseCase(deviceRepo)
	deleteUseCase := device.NewDeleteUseCase(deviceRepo)
	findUseCase := device.NewFindUseCase(deviceRepo)
	listDeviceUC := device.NewListUseCase(deviceRepo)
	updateStatusUseCase := device.NewUpdateStatusUseCase(deviceRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(registerUseCase)
	deviceHandler := handlers.NewDeviceHandler(createDeviceUC, listDeviceUC, findUseCase, updateStatusUseCase, deleteUseCase)

	return &Container{
		Logger:        logger,
		Config:        cfg,
		DB:            db,
		UserHandler:   userHandler,
		DeviceHandler: deviceHandler,
		HealthUseCase: healthUseCase,
		HealthHandler: healthHandler,
	}, nil
}

func loggerSetup(cfg *config.Config) *logrus.Logger {
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

// waitForDatabase ожидает подключения к БД с повторными попытками
func waitForDatabase(ctx context.Context, dbURL string, logger *logrus.Logger) (*pgxpool.Pool, error) {
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
