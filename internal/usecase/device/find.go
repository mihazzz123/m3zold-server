package device

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/device"
)

type FindUseCase struct {
	Repo device.Repository
}

func NewFindUseCase(repo device.Repository) *FindUseCase {
	return &FindUseCase{Repo: repo}
}

func (uc *FindUseCase) Execute(ctx context.Context, id int) (*device.Device, error) {
	return uc.Repo.FindByID(ctx, id)
}
