package container

import (
	"github.com/mihazzz123/m3zold-server/internal/delivery/http"
	database "github.com/mihazzz123/m3zold-server/internal/infrastructure"
	"github.com/mihazzz123/m3zold-server/internal/infrastructure/postgres"
	"github.com/mihazzz123/m3zold-server/internal/usecase"
	"github.com/mihazzz123/m3zold-server/internal/usecase/device"
	"github.com/mihazzz123/m3zold-server/internal/usecase/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	// Health
	HealthChecker *database.HealthChecker
	HealthUseCase *usecase.HealthUseCase
	HealthHandler *http.HealthHandler

	UserHandler   *http.UserHandler
	DeviceHandler *http.DeviceHandler
}

func New(db *pgxpool.Pool) *Container {
	// Health dependencies
	healthChecker := database.NewHealthChecker(db)
	healthUseCase := usecase.NewHealthUseCase(healthChecker)
	healthHandler := http.NewHealthHandler(healthUseCase)

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
		HealthChecker: healthChecker,
		HealthUseCase: healthUseCase,
		HealthHandler: healthHandler,
	}
}
