package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mihazzz123/m3zold-server/internal/container"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http/handlers"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http/middleware"
)

type Router struct {
	UserHandler    *handlers.UserHandler
	DeviceHandler  *handlers.DeviceHandler
	AuthMiddleware gin.HandlerFunc // позже заменим на реальную JWT-мидлвару
}

func NewRouter(di *container.Container) *gin.Engine {
	rateLimiter := middleware.NewRateLimiter()
	r := gin.Default()
	// Ограничение частоты запросов (реализуйте отдельно)
	r.Use(rateLimiter.Middleware())

	// Healthcheck
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Public routes
	r.GET("/health", di.HealthHandler.HealthCheck)
	r.GET("/ready", di.HealthHandler.ReadyCheck)

	// Авторизация
	r.POST("/auth/register", di.UserHandler.Register)
	// r.POST("/login", userHandler.Login)

	// Защищённые маршруты
	auth := r.Group("/", middleware.Auth(di.Config))
	{
		auth.POST("/devices", di.DeviceHandler.Create)
		auth.GET("/devices", di.DeviceHandler.List)
		auth.GET("/devices/:id", di.DeviceHandler.Find)
		auth.PATCH("/devices/:id/status", di.DeviceHandler.UpdateStatus)
		auth.DELETE("/devices/:id", di.DeviceHandler.Delete)
	}

	return r
}
