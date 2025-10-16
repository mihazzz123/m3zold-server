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
	// Global middleware
	r.Use(middleware.ContextMiddleware(di.Logger))

	// Public routes
	r.GET("/health", di.HealthHandler.HealthCheck)
	r.GET("/ready", di.HealthHandler.ReadyCheck)

	// Auth routes (public)
	auth := r.Group("/auth")
	{
		auth.POST("/register", di.AuthHandler.Register)
		auth.POST("/login", di.AuthHandler.Login)
		auth.POST("/refresh", di.AuthHandler.RefreshToken)
		auth.POST("/verify-email", di.VerificationEmailHandler.VerifyEmail)
		auth.POST("/resend-verification", di.VerificationEmailHandler.ResendVerification)
	}

	// Protected routes (require authentication)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Auth protected routes
		protected.POST("/auth/logout", di.AuthHandler.Logout)
		protected.POST("/auth/change-password", di.AuthHandler.ChangePassword)

		// User protected routes
		protected.GET("/users/profile", di.UserHandler.GetProfile)
		protected.PUT("/users/profile", di.UserHandler.UpdateProfile)

		// Device protected routes
		protected.POST("/devices", di.DeviceHandler.Create)
		protected.GET("/devices", di.DeviceHandler.List)
		protected.GET("/devices/:id", di.DeviceHandler.Find)
		protected.PATCH("/devices/:id/status", di.DeviceHandler.UpdateStatus)
		protected.DELETE("/devices/:id", di.DeviceHandler.Delete)

		// Email protected routes
		protected.POST("/email/welcome", di.VerificationEmailHandler.SendWelcomeEmail)
	}

	return r
}
