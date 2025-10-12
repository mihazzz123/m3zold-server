package container

import (
	"github.com/mihazzz123/m3zold-server/internal/config"
	"github.com/mihazzz123/m3zold-server/internal/delivery/http"
	"github.com/mihazzz123/m3zold-server/internal/infrastructure"
	"github.com/mihazzz123/m3zold-server/internal/infrastructure/postgres"
	"github.com/mihazzz123/m3zold-server/internal/usecase/device"
	"github.com/mihazzz123/m3zold-server/internal/usecase/health"
	"github.com/mihazzz123/m3zold-server/internal/usecase/user"
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

	AuthService *infrastructure.AuthService
}

func New(db *pgxpool.Pool, cfg *config.Config) *Container {
	// Health dependencies
	repoHealth := postgres.NewRepoHealth(db)
	checkUC := health.NewCheckUseCase(repoHealth)
	checkTablesUC := health.NewCheckTablesUseCase(repoHealth)
	getDatabaseInfoUC := health.NewGetDatabaseInfoUseCase(repoHealth)
	monitorDDUC := health.NewMonitorDBUseCase(repoHealth)
	testDBConnectionUC := health.NewTestDBConnectionUseCase(repoHealth)
	testDBConnectionHandler := http.NewTestDBConnectionHandler(testDBConnectionUC)

	healthHandler := http.NewHealthHandler(checkUC, checkTablesUC, getDatabaseInfoUC, monitorDDUC, testDBConnectionUC)

	// Auth dependencies
	authService := infrastructure.NewAuthService()

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
		HealthHandler:           healthHandler,
		TestDBConnectionHandler: testDBConnectionHandler,
		UserHandler:             userHandler,
		DeviceHandler:           deviceHandler,
		AuthService:             authService,
	}
}
