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

func (r *CachedRepository) GetTasks(ctx context.Context, limit, offset *int, userID *uuid.UUID) ([]domain.Task, error) {
	key, field := tasksListKey(userID), tasksListField(limit, offset)
	data, cacheErr := r.pool.HGet(ctx, key, field).Bytes()
	if cacheErr == nil {
		var models []TaskModel
		if err := json.Unmarshal(data, &models); err == nil && models != nil {
			return modelsToDomains(models), nil
		} else {
			if err == nil {
				err = errors.New("cached task list is null")
			}

			r.logger.Error("decode cached task list", zap.String("key", key), zap.String("field", field), zap.Error(err))
		}
	} else if !errors.Is(cacheErr, core_redis_pool.NotFound) {
		r.logger.Error("read task list from Redis cache", zap.String("key", key), zap.String("field", field), zap.Error(cacheErr))
	}

	tasks, err := r.mainRepository.GetTasks(ctx, limit, offset, userID)
	if err != nil {
		return nil, err
	}

	data, err = json.Marshal(domainsToModels(tasks))
	if err != nil {
		r.logger.Error("encode task list for Redis cache", zap.String("key", key), zap.Error(err))
		return tasks, nil
	}
	if err := r.pool.HSet(ctx, key, field, data).Err(); err != nil {
		r.logger.Error("write task list to Redis cache", zap.String("key", key), zap.String("field", field), zap.Error(err))
		return tasks, nil
	}

	expired, err := r.pool.Expire(ctx, key, r.pool.TTL()).Result()
	if err != nil || !expired {
		if err != nil {
			r.logger.Error("set task list cache TTL", zap.String("key", key), zap.Error(err))
		} else {
			r.logger.Error("task list cache key was absent when setting TTL", zap.String("key", key))
		}
		if err := r.pool.Del(ctx, key).Err(); err != nil {
			r.logger.Error("remove task list without TTL", zap.String("key", key), zap.Error(err))
		}
	}

	return tasks, nil
}
