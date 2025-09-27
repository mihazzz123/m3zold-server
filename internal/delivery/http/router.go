package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router struct {
	UserHandler    *UserHandler
	DeviceHandler  *DeviceHandler
	AuthMiddleware gin.HandlerFunc // позже заменим на реальную JWT-мидлвару
}

func NewRouter(
	userHandler *UserHandler,
	deviceHandler *DeviceHandler,
	authMiddleware gin.HandlerFunc,
) *gin.Engine {
	r := gin.Default()

	// Healthcheck
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Авторизация
	r.POST("/auth/register", userHandler.Register)
	// r.POST("/auth/login", userHandler.Login) // позже добавим

	// Защищённые маршруты
	auth := r.Group("/", authMiddleware)
	{
		auth.POST("/devices", deviceHandler.Create)
		auth.GET("/devices", deviceHandler.List)
		auth.GET("/devices/:id", deviceHandler.Find)
		auth.PATCH("/devices/:id/status", deviceHandler.UpdateStatus)
		auth.DELETE("/devices/:id", deviceHandler.Delete)
	}

	return r
}
