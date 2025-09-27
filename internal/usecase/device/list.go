package device

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/device"
)

type ListUseCase struct {
	Repo device.Repository
}

func NewListUseCase(repo device.Repository) *ListUseCase {
	return &ListUseCase{Repo: repo}
}

func (uc *ListUseCase) Execute(ctx context.Context, userID int) ([]device.Device, error) {
	return uc.Repo.ListByUser(ctx, userID)
}
