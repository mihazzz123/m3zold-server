package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mihazzz123/m3zold-server/internal/container"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http"
	"github.com/mihazzz123/m3zold-server/migrations"
)

func main() {
	// Создаем контекст для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Initialize DI container
	container, err := container.New(ctx)
	if err != nil {
		panic(err)
	}
	defer container.Close() // ✅ Теперь Close() реализован

	container.Logger.Info("🚀 Application starting...")

	// Создаем отдельный контекст для миграций
	migrateCtx, migrateCancel := context.WithTimeout(ctx, 30*time.Second)
	defer migrateCancel()

	// Run migrations
	if err := migrations.Migrate(migrateCtx, container.DB); err != nil {
		container.Logger.Fatal("Database migrations failed:", err)
	}

	container.Logger.Info("✅ Database migrations completed")

	// Initialize router
	router := http.NewRouter(container)

	// Start server in goroutine
	serverAddr := ":" + string(container.Config.App.Port)
	go func() {
		container.Logger.Infof("🌐 Server starting on %s", serverAddr)
		if err := router.Run(serverAddr); err != nil {
			container.Logger.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	container.Logger.Info("🛑 Shutdown signal received")

	// Graceful shutdown с таймаутом
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Здесь можно добавить дополнительную логику graceful shutdown
	// Например: закрытие HTTP сервера, ожидание завершения запросов и т.д.

	container.Logger.Info("👋 Application stopped gracefully")
}
