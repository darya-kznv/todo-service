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
	cacheCtx, cancel := cacheRequestContext(ctx)
	data, cacheErr := r.pool.HGet(cacheCtx, key, field).Bytes()
	cancel()
	if cacheErr == nil {
		var models []TaskModel
		if err := json.Unmarshal(data, &models); err == nil && models != nil {
			return modelsToDomains(models), nil
		} else {
			if err == nil {
				err = errors.New("cached task list is null")
			}
			r.logCacheError("decode cached task list", err, zap.String("key", key), zap.String("field", field))
		}
	} else if !errors.Is(cacheErr, core_redis_pool.NotFound) {
		r.logCacheError("read task list from Redis cache", cacheErr, zap.String("key", key), zap.String("field", field))
	}

	tasks, err := r.mainRepository.GetTasks(ctx, limit, offset, userID)
	if err != nil {
		return nil, err
	}
	data, err = json.Marshal(domainsToModels(tasks))
	if err != nil {
		r.logCacheError("encode task list for Redis cache", err, zap.String("key", key))
		return tasks, nil
	}
	cacheCtx, cancel = cacheRequestContext(ctx)
	defer cancel()
	if err := r.pool.HSet(cacheCtx, key, field, data).Err(); err != nil {
		r.logCacheError("write task list to Redis cache", err, zap.String("key", key), zap.String("field", field))
		return tasks, nil
	}
	// Hash fields share one TTL; if setting it fails, remove the hash to avoid an unbounded stale entry.
	expired, err := r.pool.Expire(cacheCtx, key, r.pool.TTL()).Result()
	if err != nil || !expired {
		if err != nil {
			r.logCacheError("set task list cache TTL", err, zap.String("key", key))
		} else {
			r.logger.Error("task list cache key was absent when setting TTL", zap.String("key", key))
		}
		cleanupCtx, cleanupCancel := cacheWriteContext(ctx)
		defer cleanupCancel()
		r.logCacheError("remove task list without TTL", r.pool.Del(cleanupCtx, key).Err(), zap.String("key", key))
	}
	return tasks, nil
}
