package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthRepository репозиторий для проверки здоровья
type HealthRepository struct {
	db *pgxpool.Pool
}

// NewHealthRepo создает новый HealthRepository
func NewHealthRepo(db *pgxpool.Pool) *HealthRepository {
	return &HealthRepository{db: db}
}

// CheckDatabase проверяет состояние базы данных
func (r *HealthRepository) CheckDatabase(ctx context.Context) (map[string]interface{}, error) {
	info := make(map[string]interface{})

	// Проверяем подключение
	if err := r.db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Получаем версию PostgreSQL
	var version string
	err := r.db.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		return nil, fmt.Errorf("failed to get database version: %w", err)
	}
	info["version"] = version

	// Получаем активные подключения
	var activeConnections int32
	err = r.db.QueryRow(ctx,
		"SELECT count(*) FROM pg_stat_activity WHERE state = 'active'").Scan(&activeConnections)
	if err == nil {
		info["active_connections"] = activeConnections
	}

	// Статистика пула
	stats := r.db.Stat()
	info["max_connections"] = stats.MaxConns()
	info["total_connections"] = stats.TotalConns()
	info["idle_connections"] = stats.IdleConns()

	return info, nil
}

// PingDatabase проверяет доступность базы данных
func (r *HealthRepository) PingDatabase(ctx context.Context) error {
	return r.db.Ping(ctx)
}
