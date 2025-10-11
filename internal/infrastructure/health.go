package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthChecker struct {
	pool *pgxpool.Pool
}

func NewHealthChecker(pool *pgxpool.Pool) *HealthChecker {
	return &HealthChecker{pool: pool}
}

// TestConnection проверяет подключение к базе данных
func (h *HealthChecker) TestConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := h.pool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	var result int
	err := h.pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("test query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected test query result: %d", result)
	}

	return nil
}

// CheckTables проверяет существование необходимых таблиц
func (h *HealthChecker) CheckTables(ctx context.Context, requiredTables []string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for _, table := range requiredTables {
		var exists bool
		query := `SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)`

		err := h.pool.QueryRow(ctx, query, table).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}

		if !exists {
			return fmt.Errorf("required table %s does not exist", table)
		}
	}

	return nil
}

// GetDatabaseInfo возвращает информацию о БД
func (h *HealthChecker) GetDatabaseInfo(ctx context.Context) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	info := make(map[string]interface{})

	// Версия PostgreSQL
	var version string
	err := h.pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err == nil {
		info["version"] = version
	}

	// Активные подключения
	var activeConnections int32
	err = h.pool.QueryRow(ctx,
		"SELECT count(*) FROM pg_stat_activity WHERE state = 'active'").Scan(&activeConnections)
	if err == nil {
		info["active_connections"] = activeConnections
	}

	// Статистика пула
	stats := h.pool.Stat()
	info["pool_max_connections"] = stats.MaxConns()
	info["pool_total_connections"] = stats.TotalConns()
	info["pool_idle_connections"] = stats.IdleConns()

	return info, nil
}
