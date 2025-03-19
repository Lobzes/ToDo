package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	
	
	"todo-api/config"
	"todo-api/database"
	"todo-api/routes"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация соединения с базой данных
	db, err := database.Connect(cfg.DBConnString)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создание экземпляра Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Базовый обработчик ошибок
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Добавление middleware для логирования запросов
	app.Use(logger.New())

	// Настройка маршрутов
	routes.Setup(app, db)

	// Запуск сервера
	log.Printf("Сервер запущен на порту %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}