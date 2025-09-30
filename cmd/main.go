package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mihazzz123/m3zold-server/internal/container"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	dbUrl := os.Getenv("DB_URL")

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("DB error:", err)
	}

	c := container.New(pool)

	r := http.NewRouter(
		c.UserHandler,
		c.DeviceHandler,
		dummyAuthMiddleware(), // временно
	)

	r.Run(":8080")
}

func dummyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Временно: берём user_id из заголовка
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Неавторизован"})
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}
