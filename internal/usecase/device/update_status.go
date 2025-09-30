package device

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/device"
)

type UpdateStatusUseCase struct {
	Repo device.Repository
}

func NewUpdateStatusUseCase(repo device.Repository) *UpdateStatusUseCase {
	return &UpdateStatusUseCase{Repo: repo}
}

func (uc *UpdateStatusUseCase) Execute(ctx context.Context, id int, status string) error {
	return uc.Repo.UpdateStatus(ctx, id, status)
}
