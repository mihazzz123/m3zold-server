package usecase

import (
	"context"
	"time"

	database "github.com/mihazzz123/m3zold-server/internal/infrastructure"
)

type HealthUseCase struct {
	healthChecker *database.HealthChecker
}

func NewHealthUseCase(healthChecker *database.HealthChecker) *HealthUseCase {
	return &HealthUseCase{
		healthChecker: healthChecker,
	}
}

type HealthStatus struct {
	Status       string                 `json:"status"`
	Timestamp    time.Time              `json:"timestamp"`
	Service      string                 `json:"service"`
	Database     string                 `json:"database"`
	DatabaseInfo map[string]interface{} `json:"database_info,omitempty"`
	Error        string                 `json:"error,omitempty"`
}

func (uc *HealthUseCase) CheckHealth(ctx context.Context) HealthStatus {
	response := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
		Service:   "m3zold-server",
		Database:  "healthy",
	}

	// Проверяем подключение к БД
	if err := uc.healthChecker.TestConnection(ctx); err != nil {
		response.Status = "degraded"
		response.Database = "unhealthy"
		response.Error = err.Error()
		return response
	}

	// Получаем дополнительную информацию о БД
	dbInfo, err := uc.healthChecker.GetDatabaseInfo(ctx)
	if err == nil {
		response.DatabaseInfo = dbInfo
	}

	return response
}

// MonitorHealth для фонового мониторинга
func (uc *HealthUseCase) MonitorHealth(ctx context.Context) error {
	return uc.healthChecker.TestConnection(ctx)
}
