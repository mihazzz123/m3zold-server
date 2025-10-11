package health

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/health"
)

type TestDBConnectionUseCase struct {
	Repo health.Repository
}

func NewTestDBConnectionUseCase(repo health.Repository) *TestDBConnectionUseCase {
	return &TestDBConnectionUseCase{
		Repo: repo,
	}
}

// TestConnection для фонового мониторинга
func (uc *TestDBConnectionUseCase) TestDBConnection(ctx context.Context) error {
	return uc.Repo.TestDBConnection(ctx)
}
