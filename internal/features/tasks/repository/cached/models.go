package tasks_cached_repository

import (
	"time"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

type TaskModel struct {
	ID           uuid.UUID  `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
	AuthorUserID uuid.UUID  `json:"author_user_id"`
}

func domainToModel(task domain.Task) TaskModel {
	return TaskModel{
		ID: task.ID, Version: task.Version, Title: task.Title,
		Description: task.Description, Completed: task.Completed,
		CompletedAt: task.CompletedAt, CreatedAt: task.CreatedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func modelToDomain(task TaskModel) domain.Task {
	return domain.NewTask(task.ID, task.Version, task.Title, task.Description,
		task.Completed, task.CompletedAt, task.CreatedAt, task.AuthorUserID)
}

func domainsToModels(tasks []domain.Task) []TaskModel {
	models := make([]TaskModel, len(tasks))
	for i, task := range tasks {
		models[i] = domainToModel(task)
	}
	return models
}

func modelsToDomains(models []TaskModel) []domain.Task {
	tasks := make([]domain.Task, len(models))
	for i, model := range models {
		tasks[i] = modelToDomain(model)
	}
	return tasks
}
