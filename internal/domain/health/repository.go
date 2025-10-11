package health

import (
	"context"
)

type Repository interface {
	Check(ctx context.Context) HealthStatus
	TestDBConnection(ctx context.Context) error
	CheckTables(ctx context.Context, requiredTables []string) error
	GetDatabaseInfo(ctx context.Context) (map[string]interface{}, error)
	MonitorDB(ctx context.Context)
}
