package tasks_postgres_repository

import (
	"time"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

type TaskModel struct {
	ID           uuid.UUID
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CompletedAt  *time.Time
	CreatedAt    time.Time
	AuthorUserID uuid.UUID
}

func taskDomainsFromModels(tasks []TaskModel) []domain.Task {
	taskDomains := make([]domain.Task, len(tasks))

	for i, task := range tasks {
		taskDomains[i] = domain.NewTask(
			task.ID,
			task.Version,
			task.Title,
			task.Description,
			task.Completed,
			task.CompletedAt,
			task.CreatedAt,
			task.AuthorUserID,
		)
	}

	return taskDomains
}
