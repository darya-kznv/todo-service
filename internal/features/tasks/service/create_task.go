package tasks_service

import (
	"context"
	"fmt"

	"github.com/darya-kznv/todo-service/internal/core/domain"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}

	taskWithUUID := domain.NewTaskWithUUID(task.Title, task.Description, task.AuthorUserID)

	createdTask, err := s.tasksRepository.CreateTask(ctx, taskWithUUID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return createdTask, nil
}
