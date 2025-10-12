package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mihazzz123/m3zold-server/internal/config"
	"github.com/mihazzz123/m3zold-server/internal/container"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http"
	"github.com/mihazzz123/m3zold-server/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.New()
	ctx := context.Background()
	logger := loggerSetup(cfg)

	// Логируем информацию о подключении к БД (без пароля)
	logger.WithFields(logrus.Fields{
		"host": cfg.Database.Host,
		"port": cfg.Database.Port,
		"user": cfg.Database.User,
		"db":   cfg.Database.DBName,
	}).Info("🔗 Initializing database connection")

	// Сначала проверяем и ждем подключение к БД
	pool, err := waitForDatabase(ctx, cfg.Database.Url, logger)
	if err != nil {
		logger.Fatal("DB connection failed:", err)
	}
	defer pool.Close()

	logger.Info("✅ Database connection established")

	// Затем выполняем миграции
	if err := migrations.Migrate(ctx, pool); err != nil {
		logger.Fatal("Database migrations failed:", err)
	}

	logger.Info("✅ Database migrations completed")

	c := container.New(pool, cfg)

	// Запускаем фоновый мониторинг здоровья БД
	go c.HealthHandler.MonitorDB(cfg)

	r := http.NewRouter(
		ctx,
		cfg,
		c.UserHandler,
		c.DeviceHandler,
		c.HealthHandler,
		c.AuthService,
	)

	serverAddr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Infof("🚀 Server starting on %s", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		logger.Fatal("Server failed to start:", err)
	}
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
