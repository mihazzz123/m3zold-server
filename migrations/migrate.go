package migrations

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// Migrate выполняет миграции из SQL файлов
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	logrus.Info("🔄 Starting database migrations...")

	// Проверяем существование папки migrations
	migrationsPath := "./migrations"
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", migrationsPath)
	}

	// Получаем список SQL файлов
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	if len(files) == 0 {
		logrus.Info("⚠️  No migration files found")
		return nil
	}

	// Сортируем файлы по имени для последовательного выполнения
	sort.Strings(files)

	// В development среде очищаем старые миграции
	if os.Getenv("APP_ENV") == "development" {
		logrus.Info("🧹 Development environment - clearing old migrations...")
		if err := clearMigrations(ctx, pool); err != nil {
			logrus.Infof("⚠️  Failed to clear old migrations: %v", err)
		}
	}

	for i, file := range files {
		migrationName := filepath.Base(file)

		// Проверяем, была ли уже выполнена эта миграция
		alreadyExecuted, err := isMigrationExecuted(ctx, pool, migrationName)
		if err != nil {
			logrus.Errorf("failed to check migration status: %s", err)
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if alreadyExecuted {
			logrus.Infof("⏭️  Migration already executed: %s", migrationName)
			continue
		}

		logrus.Infof("📁 Processing migration: %s", migrationName)

		// Читаем содержимое файла
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Выполняем миграцию в транзакции
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", file, err)
		}

		// Разделяем SQL на отдельные запросы
		queries := splitSQL(string(content))
		for j, query := range queries {
			if strings.TrimSpace(query) == "" {
				continue
			}

			// Выполняем каждый запрос отдельно
			if _, err := tx.Exec(ctx, query); err != nil {
				tx.Rollback(ctx)
				return fmt.Errorf("failed to execute query %d in migration %s: %w\nQuery: %s", j+1, file, err, query)
			}
		}

		// Отмечаем миграцию как выполненную
		if err := markMigrationAsExecuted(ctx, tx, migrationName); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to mark migration as executed: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction for %s: %w", file, err)
		}

		logrus.Infof("✅ Migration %d/%d completed: %s", i+1, len(files), migrationName)
	}

	log.Println("✅ All migrations completed successfully")
	return nil
}

func clearMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	// Проверяем что пул не закрыт и валиден
	if pool == nil {
		return fmt.Errorf("database pool is nil")
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("database pool is not connected: %w", err)
	}

	// Сначала удаляем таблицы в правильном порядке (из-за foreign keys)
	tables := []string{
		"m3zold_schema.verification_tokens",
		"m3zold_schema.devices",
		"m3zold_schema.users",
		"m3zold_schema.schema_migrations",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)
		_, err := pool.Exec(ctx, query)
		if err != nil {
			// Игнорируем ошибки "table does not exist"
			if !strings.Contains(err.Error(), "does not exist") {
				log.Printf("⚠️  Failed to drop table %s: %v", table, err)
			}
		} else {
			log.Printf("🗑️  Dropped table: %s", table)
		}
	}

	return nil
}

// splitSQL разделяет SQL файл на отдельные запросы
func splitSQL(sql string) []string {
	queries := strings.Split(sql, ";")
	var result []string

	for _, query := range queries {
		trimmed := strings.TrimSpace(query)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// isMigrationExecuted проверяет, была ли миграция уже выполнена
func isMigrationExecuted(ctx context.Context, pool *pgxpool.Pool, migrationName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM m3zold_schema.schema_migrations WHERE name = $1)`
	var exists bool
	err := pool.QueryRow(ctx, query, migrationName).Scan(&exists)
	if err != nil {
		// Если таблицы schema_migrations еще нет, считаем что миграция не выполнена
		if strings.Contains(err.Error(), "does not exist") {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

// markMigrationAsExecuted отмечает миграцию как выполненную
func markMigrationAsExecuted(ctx context.Context, tx pgx.Tx, migrationName string) error {
	query := `INSERT INTO m3zold_schema.schema_migrations (name) VALUES ($1)`
	_, err := tx.Exec(ctx, query, migrationName)
	return err
}
