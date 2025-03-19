package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"

	"todo-api/database"
	"todo-api/models"
)

// TaskHandler представляет обработчик задач
type TaskHandler struct {
	DB *database.DB
}

// NewTaskHandler создает новый обработчик задач
func NewTaskHandler(db *database.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// CreateTask обрабатывает запрос на создание новой задачи
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	// Парсим запрос
	task := new(models.CreateTaskRequest)
	if err := c.BodyParser(task); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Некорректный формат данных")
	}

	// Проверяем обязательные поля
	if task.Title == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Поле title обязательно")
	}

	// Если статус не указан, устанавливаем по умолчанию
	if task.Status == "" {
		task.Status = models.StatusNew
	} else if task.Status != models.StatusNew && task.Status != models.StatusInProgress && task.Status != models.StatusDone {
		return fiber.NewError(fiber.StatusBadRequest, "Недопустимый статус. Разрешены: new, in_progress, done")
	}

	// Создаем задачу в базе данных
	var newTask models.Task
	err := h.DB.Pool.QueryRow(
		context.Background(),
		`INSERT INTO tasks (title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		RETURNING id, title, description, status, created_at, updated_at`,
		task.Title, task.Description, task.Status,
	).Scan(&newTask.ID, &newTask.Title, &newTask.Description, &newTask.Status, &newTask.CreatedAt, &newTask.UpdatedAt)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при создании задачи: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newTask)
}

// GetAllTasks обрабатывает запрос на получение всех задач
func (h *TaskHandler) GetAllTasks(c *fiber.Ctx) error {
	// Получаем задачи из базы данных
	rows, err := h.DB.Pool.Query(
		context.Background(),
		`SELECT id, title, description, status, created_at, updated_at 
		FROM tasks ORDER BY id`,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при получении задач: "+err.Error())
	}
	defer rows.Close()

	// Преобразуем результат в срез задач
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status,
			&task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при сканировании строки: "+err.Error())
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при обработке результатов: "+err.Error())
	}

	return c.JSON(tasks)
}

// UpdateTask обрабатывает запрос на обновление задачи
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	// Получаем ID задачи из параметров URL
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Некорректный ID задачи")
	}

	// Парсим запрос
	updateReq := new(models.UpdateTaskRequest)
	if err := c.BodyParser(updateReq); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Некорректный формат данных")
	}

	// Проверяем, существует ли задача
	var existingTask models.Task
	err = h.DB.Pool.QueryRow(
		context.Background(),
		`SELECT id, title, description, status FROM tasks WHERE id = $1`,
		id,
	).Scan(&existingTask.ID, &existingTask.Title, &existingTask.Description, &existingTask.Status)

	if err == pgx.ErrNoRows {
		return fiber.NewError(fiber.StatusNotFound, "Задача не найдена")
	} else if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при поиске задачи: "+err.Error())
	}

	// Обновляем только предоставленные поля
	title := existingTask.Title
	if updateReq.Title != "" {
		title = updateReq.Title
	}

	description := existingTask.Description
	if updateReq.Description != "" {
		description = updateReq.Description
	}

	status := existingTask.Status
	if updateReq.Status != "" {
		if updateReq.Status != models.StatusNew && updateReq.Status != models.StatusInProgress && updateReq.Status != models.StatusDone {
			return fiber.NewError(fiber.StatusBadRequest, "Недопустимый статус. Разрешены: new, in_progress, done")
		}
		status = updateReq.Status
	}

	// Обновляем задачу в базе данных
	var updatedTask models.Task
	err = h.DB.Pool.QueryRow(
		context.Background(),
		`UPDATE tasks 
		SET title = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, title, description, status, created_at, updated_at`,
		title, description, status, time.Now(), id,
	).Scan(
		&updatedTask.ID, &updatedTask.Title, &updatedTask.Description,
		&updatedTask.Status, &updatedTask.CreatedAt, &updatedTask.UpdatedAt,
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при обновлении задачи: "+err.Error())
	}

	return c.JSON(updatedTask)
}

// DeleteTask обрабатывает запрос на удаление задачи
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	// Получаем ID задачи из параметров URL
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Некорректный ID задачи")
	}

	// Удаляем задачу из базы данных
	cmdTag, err := h.DB.Pool.Exec(
		context.Background(),
		`DELETE FROM tasks WHERE id = $1`,
		id,
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при удалении задачи: "+err.Error())
	}

	if cmdTag.RowsAffected() == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Задача не найдена")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Задача успешно удалена",
		"id":      id,
	})
}