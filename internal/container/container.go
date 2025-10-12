package container

import (
	"github.com/mihazzz123/m3zold-server/internal/config"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http"
	"github.com/mihazzz123/m3zold-server/internal/infrastructure/postgres"
	infrastructure_services "github.com/mihazzz123/m3zold-server/internal/infrastructure/services"
	"github.com/mihazzz123/m3zold-server/internal/usecase/device"
	userusecase "github.com/mihazzz123/m3zold-server/internal/usecase/user"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Config *config.Config
	Logger *logrus.Logger
	// Health
	HealthHandler           *http.HealthHandler
	TestDBConnectionHandler *http.TestDBConnectionHandler

	UserHandler   *http.UserHandler
	DeviceHandler *http.DeviceHandler
}

func New(db *pgxpool.Pool, cfg *config.Config) *Container {
	// Repositories
	userRepo := postgres.NewUserRepo(db)
	deviceRepo := postgres.NewDeviceRepo(db)

	// Services
	passwordService := infrastructure_services.NewPasswordService(0)
	idService := infrastructure_services.NewIDService()
	emailService := infrastructure_services.NewEmailService()

	// Use Cases с внедренными services
	registerUseCase := userusecase.NewRegisterUseCase(
		userRepo,
		passwordService,
		idService,
		emailService,
	)

	// Device UseCases
	createDeviceUC := device.NewCreateUseCase(deviceRepo)
	deleteUseCase := device.NewDeleteUseCase(deviceRepo)
	findUseCase := device.NewFindUseCase(deviceRepo)
	listDeviceUC := device.NewListUseCase(deviceRepo)
	updateStatusUseCase := device.NewUpdateStatusUseCase(deviceRepo)

	// Handlers
	userHandler := http.NewUserHandler(registerUseCase)
	deviceHandler := http.NewDeviceHandler(createDeviceUC, listDeviceUC, findUseCase, updateStatusUseCase, deleteUseCase)

	return &Container{
		UserHandler:   userHandler,
		DeviceHandler: deviceHandler,
	}
}
