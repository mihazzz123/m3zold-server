package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	health "github.com/mihazzz123/m3zold-server/internal/usecase"
)

type HealthHandler struct {
	healthUseCase *health.HealthUseCase
}

func NewHealthHandler(healthUseCase *health.HealthUseCase) *HealthHandler {
	return &HealthHandler{
		healthUseCase: healthUseCase,
	}
}

// HealthCheck обработчик для /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := h.healthUseCase.CheckHealth(c.Request.Context())

	if status.Status == "unhealthy" {
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}

// ReadyCheck обработчик для /ready
func (h *HealthHandler) ReadyCheck(c *gin.Context) {
	status := h.healthUseCase.CheckHealth(c.Request.Context())

	if status.Database != "healthy" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"service": "m3zold-server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"service": "m3zold-server",
	})
}
