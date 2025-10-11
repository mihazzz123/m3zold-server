package health

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/health"
)

type CheckTablesUseCase struct {
	Repo health.Repository
}

func NewCheckTablesUseCase(repo health.Repository) *CheckTablesUseCase {
	return &CheckTablesUseCase{
		Repo: repo,
	}
}

// MonitorHealth для фонового мониторинга
func (uc *CheckTablesUseCase) CheckTables(ctx context.Context, requiredTables []string) error {
	return uc.Repo.CheckTables(ctx, requiredTables)
}
