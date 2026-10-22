package tasks_cached_repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	core_redis_pool "github.com/darya-kznv/todo-service/internal/core/repository/redis/pool"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *CachedRepository) getTaskFromCache(ctx context.Context, id uuid.UUID) (domain.Task, bool) {
	data, err := r.pool.Get(ctx, taskKey(id)).Bytes()
	if err != nil {
		if !errors.Is(err, core_redis_pool.NotFound) {
			r.logger.Error("read task from Redis cache", zap.String("task_id", id.String()), zap.Error(err))
		}
		return domain.Task{}, false
	}

	var model TaskModel
	if err := json.Unmarshal(data, &model); err != nil {
		r.logger.Error("decode cached task", zap.String("task_id", id.String()), zap.Error(err))
		return domain.Task{}, false
	}
	if model.ID != id {
		r.logger.Error("cached task ID does not match key", zap.String("task_id", id.String()), zap.String("cached_id", model.ID.String()))
		return domain.Task{}, false
	}

	return modelToDomain(model), true
}

func (r *CachedRepository) cacheTask(ctx context.Context, task domain.Task) {
	data, err := json.Marshal(domainToModel(task))
	if err != nil {
		r.logger.Error("encode task for Redis cache", zap.String("task_id", task.ID.String()), zap.Error(err))
		return
	}

	if err := r.pool.Set(ctx, taskKey(task.ID), data, r.pool.TTL()).Err(); err != nil {
		r.logger.Error("write task to Redis cache", zap.String("task_id", task.ID.String()), zap.Error(err))
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
