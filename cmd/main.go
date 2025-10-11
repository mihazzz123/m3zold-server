package main

import (
	"context"
	"os"

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

	pool, err := pgxpool.New(ctx, cfg.Database.Url)
	if err != nil {
		logger.Fatal("DB error:", err)
	}
	defer pool.Close()

	migrations.Migrate(ctx, pool)

	c := container.New(pool, cfg)

	// Тест подключения к базе данных при старте
	if err := c.TestDBConnectionHandler.TestDBConnection(cfg.Database.Url); err != nil {
		logger.Fatal("DB connection test failed:", err)
	}
	// Запускаем фоновый мониторинг здоровья БД
	go c.HealthHandler.MonitorDB(cfg)

	r := http.NewRouter(
		ctx,
		cfg,
		c.UserHandler,
		c.DeviceHandler,
		c.HealthHandler,
	)

	r.Run(":8080")
}

func loggerSetup(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()
	if cfg.Logger.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else if cfg.Logger.Format == "text" {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
	if cfg.Logger.Output == "ststdoutderr" {
		logger.SetOutput(os.Stdout)
	} else if cfg.Logger.Output == "file" {
		file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer file.Close()
		if err == nil {
			logger.SetOutput(file)
		} else {
			logrus.Info("Failed to log to file, using default stderr")
		}
	}
	if cfg.Logger.Level == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else if cfg.Logger.Level == "info" {
		logger.SetLevel(logrus.InfoLevel)
	} else if cfg.Logger.Level == "error" {
		logger.SetLevel(logrus.ErrorLevel)
	} else if cfg.Logger.Level == "warn" {
		logger.SetLevel(logrus.WarnLevel)
	}
	return logger
}
