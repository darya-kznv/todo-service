package tasks_service

import (
	"context"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(
		ctx context.Context,
		domain domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		limit *int,
		offset *int,
		userID *uuid.UUID,
	) ([]domain.Task, error)

	DeleteTask(
		ctx context.Context,
		id uuid.UUID,
	) error

	GetTask(
		ctx context.Context,
		id uuid.UUID,
	) (domain.Task, error)

	PatchTask(
		ctx context.Context,
		id uuid.UUID,
		task domain.Task,
	) (domain.Task, error)
}

func NewTasksService(
	tasksRepository TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
