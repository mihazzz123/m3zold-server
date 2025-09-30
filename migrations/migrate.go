package db

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	log.Println("🔄 Запуск миграций из файлов...")

	files, err := filepath.Glob(filepath.Join("./migrations", "*.sql"))
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("❌ Ошибка чтения %s: %v\n", file, err)
			return err
		}

		_, err = pool.Exec(context.Background(), string(content))
		if err != nil {
			log.Printf("❌ Ошибка выполнения миграции %s: %v\n", file, err)
			return err
		}

		log.Printf("✅ Миграция выполнена: %s\n", filepath.Base(file))
	}

	log.Println("✅ Все миграции завершены")
	return nil
}
