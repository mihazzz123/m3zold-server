package main

import (
	"context"
	"fmt"

	"github.com/mihazzz123/m3zold-server/internal/container"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http"
	"github.com/mihazzz123/m3zold-server/migrations"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	di, err := container.New(ctx)
	if err != nil {
		di.Logger.Fatal("Container initialization failed:", err)
	}
	// Логируем информацию о подключении к БД (без пароля)
	di.Logger.WithFields(logrus.Fields{
		"host": di.Config.Database.Host,
		"port": di.Config.Database.Port,
		"user": di.Config.Database.User,
		"db":   di.Config.Database.DBName,
	}).Info("🔗 Initializing database connection")

	di.Logger.Info("✅ Database connection established")

	// Затем выполняем миграции
	if err := migrations.Migrate(ctx, di.DB); err != nil {
		di.Logger.Fatal("Database migrations failed:", err)
	}

	di.Logger.Info("✅ Database migrations completed")

	// Запускаем фоновый мониторинг здоровья БД
	go di.HealthUseCase.MonitorDB(ctx)

	r := http.NewRouter(di)

	serverAddr := fmt.Sprintf(":%d", di.Config.App.Port)
	di.Logger.Infof("🚀 Server starting on %s", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		di.Logger.Fatal("Server failed to start:", err)
	}
}
