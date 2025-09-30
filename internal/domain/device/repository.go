package device

import "context"

type Repository interface {
	Create(ctx context.Context, d *Device) error
	ListByUser(ctx context.Context, userID int) ([]Device, error)
	FindByID(ctx context.Context, id int) (*Device, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error
}
