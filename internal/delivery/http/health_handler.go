package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mihazzz123/m3zold-server/internal/config"
	"github.com/mihazzz123/m3zold-server/internal/usecase/health"
)

type HealthHandler struct {
	checkUC            *health.CheckUseCase
	CheckTablesUC      *health.CheckTablesUseCase
	GetDatabaseInfoUC  *health.GetDatabaseInfoUseCase
	MonitorDBUC        *health.MonitorDBUseCase
	TestDBConnectionUC *health.TestDBConnectionUseCase
}

func NewHealthHandler(
	checkUC *health.CheckUseCase,
	checkTablesUC *health.CheckTablesUseCase,
	getDatabaseInfoUC *health.GetDatabaseInfoUseCase,
	monitorDBUC *health.MonitorDBUseCase,
	testDBConnectionUC *health.TestDBConnectionUseCase,
) *HealthHandler {
	return &HealthHandler{
		checkUC:            checkUC,
		CheckTablesUC:      checkTablesUC,
		GetDatabaseInfoUC:  getDatabaseInfoUC,
		MonitorDBUC:        monitorDBUC,
		TestDBConnectionUC: testDBConnectionUC}
}

// HealthCheck обработчик для /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := h.checkUC.Check(c.Request.Context())

	if status.Status == "degraded" {
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}

// ReadyCheck обработчик для /ready (только когда приложение готово обслуживать трафик)
func (h *HealthHandler) ReadyCheck(c *gin.Context) {
	status := h.checkUC.Check(c.Request.Context())

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

// MonitorDB фоновая проверка здоровья БД
func (h *HealthHandler) MonitorDB(cfg *config.Config) {
	h.MonitorDBUC.MonitorDB(context.Background())
}

func (h *HealthHandler) TestDBConnection(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return err
	}
	defer pool.Close()

	return h.TestDBConnectionUC.TestDBConnection(ctx)
}
