package routes

import (
	"github.com/gofiber/fiber/v2"
	
	"todo-api/database"
	"todo-api/handlers"
)

// Setup настраивает все маршруты API
func Setup(app *fiber.App, db *database.DB) {
	// Создаем обработчик задач
	taskHandler := handlers.NewTaskHandler(db)

	// Создаем группу маршрутов для API
	api := app.Group("/")

	// Маршруты для задач
	api.Post("/tasks", taskHandler.CreateTask)
	api.Get("/tasks", taskHandler.GetAllTasks)
	api.Put("/tasks/:id", taskHandler.UpdateTask)
	api.Delete("/tasks/:id", taskHandler.DeleteTask)
}