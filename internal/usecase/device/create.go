package device

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/device"
)

type CreateInput struct {
	UserID int
	Name   string
	Type   string
}

type CreateUseCase struct {
	Repo device.Repository
}

func NewCreateUseCase(repo device.Repository) *CreateUseCase {
	return &CreateUseCase{Repo: repo}
}

func (uc *CreateUseCase) Execute(ctx context.Context, input CreateInput) error {
	d := &device.Device{
		UserID: input.UserID,
		Name:   input.Name,
		Type:   input.Type,
		Status: "offline",
	}
	return uc.Repo.Create(ctx, d)
}
