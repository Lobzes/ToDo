package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	

	"github.com/jackc/pgx/v4/pgxpool"
)

// DB представляет обертку над пулом соединений с базой данных
type DB struct {
	Pool *pgxpool.Pool
}

// Connect устанавливает соединение с базой данных PostgreSQL
func Connect(connString string) (*DB, error) {
	// Создаем пул соединений
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Проверяем подключение
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ошибка проверки подключения к базе данных: %w", err)
	}

	db := &DB{Pool: pool}

	// Применяем миграции
	if err := db.Migrate(); err != nil {
		return nil, fmt.Errorf("ошибка применения миграций: %w", err)
	}

	return db, nil
}

// Close закрывает соединение с базой данных
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Migrate применяет миграции к базе данных
func (db *DB) Migrate() error {
	// Применяем миграцию для создания таблицы задач
	migrationsDir := "database/migrations"

	// Проверяем существование директории миграций
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Если директория не существует, создаем таблицу напрямую
		return db.createTasksTable()
	}

	// Находим все SQL файлы в директории миграций
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("ошибка поиска файлов миграций: %w", err)
	}

	// Применяем миграции в порядке их номеров
	for _, file := range files {
		migrationSQL, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("ошибка чтения файла миграции %s: %w", file, err)
		}

		_, err = db.Pool.Exec(context.Background(), string(migrationSQL))
		if err != nil {
			return fmt.Errorf("ошибка выполнения миграции %s: %w", file, err)
		}

		fmt.Printf("Применена миграция: %s\n", filepath.Base(file))
	}

	return nil
}

// createTasksTable создает таблицу задач напрямую
func (db *DB) createTasksTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);
	`

	_, err := db.Pool.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы tasks: %w", err)
	}

	fmt.Println("Таблица tasks создана")
	return nil
}