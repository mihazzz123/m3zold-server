package postgres

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/device"

	"github.com/jackc/pgx/v5/pgxpool"
)

type deviceRepo struct {
	DB *pgxpool.Pool
}

func NewDeviceRepo(db *pgxpool.Pool) *deviceRepo {
	return &deviceRepo{DB: db}
}

// ✅ Create — добавление нового устройства
func (r *deviceRepo) Create(ctx context.Context, d *device.Device) error {
	_, err := r.DB.Exec(ctx,
		`INSERT INTO devices (user_id, name, type, status)
         VALUES ($1, $2, $3, $4)`,
		d.UserID, d.Name, d.Type, d.Status)
	return err
}

// ✅ ListByUser — список устройств пользователя
func (r *deviceRepo) ListByUser(ctx context.Context, userID int) ([]device.Device, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, user_id, name, type, status
         FROM devices WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []device.Device
	for rows.Next() {
		var d device.Device
		if err := rows.Scan(&d.ID, &d.UserID, &d.Name, &d.Type, &d.Status); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

// ✅ FindByID — получить устройство по ID
func (r *deviceRepo) FindByID(ctx context.Context, id int) (*device.Device, error) {
	row := r.DB.QueryRow(ctx,
		`SELECT id, user_id, name, type, status
         FROM devices WHERE id = $1`, id)

	var d device.Device
	if err := row.Scan(&d.ID, &d.UserID, &d.Name, &d.Type, &d.Status); err != nil {
		return nil, err
	}
	return &d, nil
}

// ✅ UpdateStatus — обновить статус устройства
func (r *deviceRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE devices SET status = $1 WHERE id = $2`,
		status, id)
	return err
}

// ✅ Delete — удалить устройство
func (r *deviceRepo) Delete(ctx context.Context, id int) error {
	_, err := r.DB.Exec(ctx,
		`DELETE FROM devices WHERE id = $1`, id)
	return err
}
