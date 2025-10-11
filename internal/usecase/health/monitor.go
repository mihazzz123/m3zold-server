package health

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/health"
)

type MonitorDBUseCase struct {
	Repo health.Repository
}

func NewMonitorDBUseCase(repo health.Repository) *MonitorDBUseCase {
	return &MonitorDBUseCase{
		Repo: repo,
	}
}

// MonitorHealth для фонового мониторинга
func (uc *MonitorDBUseCase) MonitorDB(ctx context.Context) {
	uc.Repo.MonitorDB(ctx)
}
