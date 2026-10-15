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
	initializedTask, err := task.Initialized(s.clock.Now())
	if err != nil {
		return domain.Task{}, fmt.Errorf("initialize task domain: %w", err)
	}

	createdTask, err := s.tasksRepository.CreateTask(ctx, initializedTask)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return createdTask, nil
}
