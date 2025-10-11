package health

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/health"
)

type CheckUseCase struct {
	Repo health.Repository
}

func NewCheckUseCase(repo health.Repository) *CheckUseCase {
	return &CheckUseCase{
		Repo: repo,
	}
}

func (uc *CheckUseCase) Check(ctx context.Context) health.HealthStatus {
	return uc.Repo.Check(ctx)
}
