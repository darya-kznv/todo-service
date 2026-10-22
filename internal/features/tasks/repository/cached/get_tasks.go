package tasks_cached_repository

import (
	"context"
	"encoding/json"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

func (r *CachedRepository) GetTasks(ctx context.Context, limit, offset *int, userID *uuid.UUID) ([]domain.Task, error) {
	key, field := tasksListKey(userID), tasksListField(limit, offset)
	if data, err := r.pool.HGet(ctx, key, field).Bytes(); err == nil {
		var models []TaskModel
		if json.Unmarshal(data, &models) == nil {
			return modelsToDomains(models), nil
		}
	}

	tasks, err := r.mainRepository.GetTasks(ctx, limit, offset, userID)
	if err != nil {
		return nil, err
	}
	if data, err := json.Marshal(domainsToModels(tasks)); err == nil {
		if r.pool.HSet(ctx, key, field, data).Err() == nil {
			// Hash fields share one TTL; an expired key is always repopulated from Postgres.
			if expired, err := r.pool.Expire(ctx, key, r.pool.TTL()).Result(); err != nil || !expired {
				_ = r.pool.Del(ctx, key).Err()
			}
		}
	}
	return tasks, nil
}
