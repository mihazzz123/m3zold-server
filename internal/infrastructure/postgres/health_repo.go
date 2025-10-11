package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mihazzz123/m3zold-server/internal/domain/health"
	"github.com/sirupsen/logrus"
)

type repoHealth struct {
	pool *pgxpool.Pool
}

func NewRepoHealth(pool *pgxpool.Pool) *repoHealth {
	return &repoHealth{pool: pool}
}

func (h *repoHealth) Check(ctx context.Context) health.HealthStatus {
	response := health.HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
		Service:   "m3zold-server",
		Database:  "healthy",
	}

	// Проверяем подключение к БД
	if err := h.TestDBConnection(ctx); err != nil {
		response.Status = "degraded"
		response.Database = "unhealthy"
		response.Error = err.Error()
		return response
	}

	// Получаем дополнительную информацию о БД
	dbInfo, err := h.GetDatabaseInfo(ctx)
	if err == nil {
		response.DatabaseInfo = dbInfo
	}

	return response
}

// TestConnection проверяет подключение к базе данных
func (h *repoHealth) TestDBConnection(ctx context.Context) error {
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
func (h *repoHealth) CheckTables(ctx context.Context, requiredTables []string) error {
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
func (h *repoHealth) GetDatabaseInfo(ctx context.Context) (map[string]interface{}, error) {
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

// MonitorDB фоновая проверка здоровья БД
func (h *repoHealth) MonitorDB(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := h.TestDBConnection(ctx)
		cancel()

		if err != nil {
			logrus.Printf("⚠️ Database health check failed: %v", err)
		} else {
			logrus.Println("✅ Database health check passed")
		}
	}
}
