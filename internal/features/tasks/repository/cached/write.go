package tasks_cached_repository

import (
	"context"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *CachedRepository) invalidateTasks(ctx context.Context, userID uuid.UUID, taskID *uuid.UUID) {
	keys := []string{tasksListKey(nil), tasksListKey(&userID)}

	if taskID != nil {
		keys = append(keys, taskKey(*taskID))
	}

	if err := r.pool.Del(ctx, keys...).Err(); err != nil {
		r.logger.Error("invalidate task cache", zap.Strings("keys", keys), zap.Error(err))
	}
}

func (r *CachedRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	created, err := r.mainRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}

	r.invalidateTasks(ctx, created.AuthorUserID, nil)
	r.cacheTask(ctx, created)

	return created, nil
}

func (r *CachedRepository) PatchTask(ctx context.Context, id uuid.UUID, task domain.Task) (domain.Task, error) {
	patched, err := r.mainRepository.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, err
	}

	r.invalidateTasks(ctx, patched.AuthorUserID, &id)
	r.cacheTask(ctx, patched)

	return patched, nil
}

func (r *CachedRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	task, ok := r.getTaskFromCache(ctx, id)
	if !ok {
		var err error

		task, err = r.mainRepository.GetTask(ctx, id)
		if err != nil {
			return err
		}
	}

	if err := r.mainRepository.DeleteTask(ctx, id); err != nil {
		return err
	}

	r.invalidateTasks(ctx, task.AuthorUserID, &id)

	return nil
}
