package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mihazzz123/m3zold-server/internal/container"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http"
	database "github.com/mihazzz123/m3zold-server/internal/infrastructure"
	"github.com/mihazzz123/m3zold-server/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	dbUrl := os.Getenv("DB_URL")

	// Тест подключения к базе данных при старте
	if err := testDBConnection(dbUrl); err != nil {
		log.Fatal("DB connection test failed:", err)
	}

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("DB error:", err)
	}
	defer pool.Close()

	c := container.New(pool)

	r := http.NewRouter(
		dummyAuthMiddleware(),
		c.UserHandler,
		c.DeviceHandler,
		c.HealthHandler,
	)

	r.Run(":8080")
}

// testDBConnection использует инфраструктурный слой
func testDBConnection(dbURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	healthChecker := database.NewHealthChecker(pool)
	return healthChecker.TestConnection(ctx)
}

// monitorDBHealth фоновая проверка здоровья БД
func monitorDBHealth(healthUseCase *usecase.HealthUseCase) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := healthUseCase.MonitorHealth(ctx)
		cancel()

		if err != nil {
			log.Printf("⚠️ Database health check failed: %v", err)
		} else {
			log.Println("✅ Database health check passed")
		}
	}
}

func dummyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Неавторизован"})
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}
