package main

import (
	"context"
	"fmt"
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
	r := http.NewRouter(container)

	// Start server in goroutine
	serverAddr := fmt.Sprintf(":%d", container.Config.App.Port)
	container.Logger.Infof("🚀 Server starting on %s", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		container.Logger.Fatal("Server failed to start:", err)
	}

	container.Logger.Info("👋 Application stopped gracefully")
}
