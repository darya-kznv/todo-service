package tasks_transport_http

import (
	"time"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

type TaskDTOResponse struct {
	ID           uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Version      int        `json:"version" example:"1"`
	Title        string     `json:"title" example:"Write API documentation"`
	Description  *string    `json:"description" example:"Document every endpoint"`
	Completed    bool       `json:"completed" example:"false"`
	CompletedAt  *time.Time `json:"completed_at" example:"2026-01-15T11:30:00Z"`
	CreatedAt    time.Time  `json:"created_at" example:"2026-01-15T10:30:00Z"`
	AuthorUserID uuid.UUID  `json:"author_user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CompletedAt:  task.CompletedAt,
		CreatedAt:    task.CreatedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func tasksDTOFromDomains(tasks []domain.Task) []TaskDTOResponse {
	tasksDTO := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		tasksDTO[i] = taskDTOFromDomain(task)
	}

	return tasksDTO
}
