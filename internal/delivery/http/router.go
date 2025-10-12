package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mihazzz123/m3zold-server/internal/config"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http/middleware"
	"github.com/mihazzz123/m3zold-server/internal/infrastructure"
)

type Router struct {
	UserHandler    *UserHandler
	DeviceHandler  *DeviceHandler
	AuthMiddleware gin.HandlerFunc // позже заменим на реальную JWT-мидлвару
}

func NewRouter(
	ctx context.Context,
	cfg *config.Config,
	userHandler *UserHandler,
	deviceHandler *DeviceHandler,
	healthHandler *HealthHandler,
	authSrv *infrastructure.AuthService,
) *gin.Engine {
	rateLimiter := middleware.NewRateLimiter()
	r := gin.Default()
	// Ограничение частоты запросов (реализуйте отдельно)
	r.Use(rateLimiter.Middleware())

	// Healthcheck
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Public routes
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/ready", healthHandler.ReadyCheck)

	// Авторизация
	r.POST("/auth/register", userHandler.Register)
	r.POST("/login", func(c *gin.Context) {
		// Пример: авторизация по userID
		userID := c.PostForm("user_id")
		userUIID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid user_id"})
			return
		}
		token, err := authSrv.GenerateToken(cfg, userUIID)
		if err != nil {
			c.JSON(500, gin.H{"error": "token generation failed"})
			return
		}
		c.JSON(200, gin.H{"token": token})
	})

	// Защищённые маршруты
	auth := r.Group("/", middleware.Auth(cfg))
	{
		auth.POST("/devices", deviceHandler.Create)
		auth.GET("/devices", deviceHandler.List)
		auth.GET("/devices/:id", deviceHandler.Find)
		auth.PATCH("/devices/:id/status", deviceHandler.UpdateStatus)
		auth.DELETE("/devices/:id", deviceHandler.Delete)
	}

	return r
}
