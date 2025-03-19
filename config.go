package config

import (
	"errors"
	"os"
)

// Config содержит конфигурационные параметры приложения
type Config struct {
	Port         string
	DBConnString string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Порт по умолчанию
	}

	dbConnString := os.Getenv("DB_CONN_STRING")
	if dbConnString == "" {
		// Строка подключения по умолчанию
		dbConnString = "postgres://postgres:postgres@localhost:5432/tododb"
	}

	// Проверка обязательных параметров
	if dbConnString == "" {
		return nil, errors.New("не указана строка подключения к базе данных")
	}

	return &Config{
		Port:         port,
		DBConnString: dbConnString,
	}, nil
}