package health

import (
	"context"
	"time"

	"github.com/mihazzz123/m3zold-server/internal/domain/health"
)

type HealthUseCase struct {
	healthRepo health.Repository
}

func NewHealthUseCase(healthRepo health.Repository) *HealthUseCase {
	return &HealthUseCase{
		healthRepo: healthRepo,
	}
}

type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Service   string                 `json:"service"`
	Database  string                 `json:"database"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

func (uc *HealthUseCase) CheckHealth(ctx context.Context) HealthStatus {
	response := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Service:   "m3zold-server",
		Database:  "healthy",
	}

	// Проверяем здоровье БД через репозиторий
	dbStatus, err := uc.healthRepo.CheckDatabase(ctx)
	if err != nil {
		response.Status = "unhealthy"
		response.Database = "unhealthy"
		response.Error = err.Error()
		return response
	}

	response.Details = map[string]interface{}{
		"database": dbStatus,
	}

	return response
}
