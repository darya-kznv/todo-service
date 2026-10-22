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
	r.logCacheError("invalidate task cache", r.pool.Del(ctx, keys...).Err(), zap.Strings("keys", keys))
}

func (r *CachedRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	created, err := r.mainRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}
	cacheCtx, cancel := cacheWriteContext(ctx)
	defer cancel()
	r.invalidateTasks(cacheCtx, created.AuthorUserID, nil)
	r.cacheTask(cacheCtx, created)
	return created, nil
}

func (r *CachedRepository) PatchTask(ctx context.Context, id uuid.UUID, task domain.Task) (domain.Task, error) {
	patched, err := r.mainRepository.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, err
	}
	cacheCtx, cancel := cacheWriteContext(ctx)
	defer cancel()
	r.invalidateTasks(cacheCtx, patched.AuthorUserID, &id)
	r.cacheTask(cacheCtx, patched)
	return patched, nil
}

func (r *CachedRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	// The author is needed to invalidate their paginated list after deletion.
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
	cacheCtx, cancel := cacheWriteContext(ctx)
	defer cancel()
	r.invalidateTasks(cacheCtx, task.AuthorUserID, &id)
	return nil
}
