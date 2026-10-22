package tasks_cached_repository

import (
	"context"
	"encoding/json"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

func (r *CachedRepository) getTaskFromCache(ctx context.Context, id uuid.UUID) (domain.Task, bool) {
	data, err := r.pool.Get(ctx, taskKey(id)).Bytes()
	if err != nil {
		return domain.Task{}, false
	}
	var model TaskModel
	if err := json.Unmarshal(data, &model); err != nil || model.ID != id {
		return domain.Task{}, false
	}
	return modelToDomain(model), true
}

func (r *CachedRepository) cacheTask(ctx context.Context, task domain.Task) {
	data, err := json.Marshal(domainToModel(task))
	if err == nil {
		_ = r.pool.Set(ctx, taskKey(task.ID), data, r.pool.TTL()).Err()
	}
}

func (r *CachedRepository) GetTask(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	if task, ok := r.getTaskFromCache(ctx, id); ok {
		return task, nil
	}
	task, err := r.mainRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}
	r.cacheTask(ctx, task)
	return task, nil
}
