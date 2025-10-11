package health

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/health"
)

type GetDatabaseInfoUseCase struct {
	Repo health.Repository
}

func NewGetDatabaseInfoUseCase(repo health.Repository) *GetDatabaseInfoUseCase {
	return &GetDatabaseInfoUseCase{
		Repo: repo,
	}
}

// MonitorHealth для фонового мониторинга
func (uc *GetDatabaseInfoUseCase) GetDatabaseInfo(ctx context.Context) (map[string]interface{}, error) {
	return uc.Repo.GetDatabaseInfo(ctx)
}
