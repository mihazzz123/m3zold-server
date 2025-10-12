package health

import "context"

type Repository interface {
	CheckDatabase(ctx context.Context) (map[string]interface{}, error)
	PingDatabase(ctx context.Context) error
}
