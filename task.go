package models

import "time"

// Статусы задач
const (
	StatusNew        = "new"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

// Task представляет модель задачи
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTaskRequest представляет запрос на создание задачи
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status,omitempty"`
}

// UpdateTaskRequest представляет запрос на обновление задачи
type UpdateTaskRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}