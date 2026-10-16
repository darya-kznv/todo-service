package tasks_service

import (
	"context"

	core_clock "github.com/darya-kznv/todo-service/internal/core/clock"
	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

type TasksService struct {
	tasksRepository TasksRepository
	clock           core_clock.Clock
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
}

func NewTasksService(
	tasksRepository TasksRepository,
	clock core_clock.Clock,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
		clock:           clock,
	}
}
