package tasks_service

import (
	"context"
	"fmt"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	id uuid.UUID,
	patch domain.TaskPatch,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if err := task.ApplyPatch(patch, s.clock.Now()); err != nil {
		return domain.Task{}, fmt.Errorf("apply task patch: %w", err)
	}

	patchedTask, err := s.tasksRepository.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}

	return patchedTask, nil
}
