package tasks_service

import (
	"context"
	"fmt"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	id uuid.UUID,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	return task, nil
}
