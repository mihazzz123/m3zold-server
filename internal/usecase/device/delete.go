package device

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/device"
)

type DeleteUseCase struct {
	Repo device.Repository
}

func NewDeleteUseCase(repo device.Repository) *DeleteUseCase {
	return &DeleteUseCase{Repo: repo}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, id int) error {
	return uc.Repo.Delete(ctx, id)
}
