package container

import (
	"github.com/mihazzz123/m3zold-server/internal/delivery/http"
	"github.com/mihazzz123/m3zold-server/internal/infrastructure/postgres"
	"github.com/mihazzz123/m3zold-server/internal/usecase/device"
	"github.com/mihazzz123/m3zold-server/internal/usecase/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	UserHandler   *http.UserHandler
	DeviceHandler *http.DeviceHandler
}

func New(db *pgxpool.Pool) *Container {
	// Repositories
	userRepo := postgres.NewUserRepo(db)
	deviceRepo := postgres.NewDeviceRepo(db)

	// User UseCases
	registerUC := user.NewRegisterUseCase(userRepo)
	// Device UseCases
	createDeviceUC := device.NewCreateUseCase(deviceRepo)
	deleteUseCase := device.NewDeleteUseCase(deviceRepo)
	findUseCase := device.NewFindUseCase(deviceRepo)
	listDeviceUC := device.NewListUseCase(deviceRepo)
	updateStatusUseCase := device.NewUpdateStatusUseCase(deviceRepo)

	// Handlers
	userHandler := http.NewUserHandler(registerUC)
	deviceHandler := http.NewDeviceHandler(createDeviceUC, listDeviceUC, findUseCase, updateStatusUseCase, deleteUseCase)

	return &Container{
		UserHandler:   userHandler,
		DeviceHandler: deviceHandler,
	}
}
