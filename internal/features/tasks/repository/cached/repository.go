package tasks_cached_repository

import (
	"context"
	"fmt"
	"time"

	core_redis_pool "github.com/darya-kznv/todo-service/internal/core/repository/redis/pool"
	tasks_service "github.com/darya-kznv/todo-service/internal/features/tasks/service"
	"github.com/google/uuid"
)

type CachedRepository struct {
	pool           core_redis_pool.Pool
	mainRepository tasks_service.TasksRepository
}

var _ tasks_service.TasksRepository = (*CachedRepository)(nil)

func NewCachedRepository(pool core_redis_pool.Pool, mainRepository tasks_service.TasksRepository) *CachedRepository {
	return &CachedRepository{pool: pool, mainRepository: mainRepository}
}

func taskKey(id uuid.UUID) string { return "task:" + id.String() }

func tasksListKey(userID *uuid.UUID) string {
	if userID == nil {
		return "tasks:all"
	}
	return "tasks:" + userID.String()
}

func tasksListField(limit, offset *int) string {
	format := func(value *int) string {
		if value == nil {
			return "nil"
		}
		return fmt.Sprint(*value)
	}
	return format(limit) + ":" + format(offset)
}

// A completed database write must still invalidate cache if the request was cancelled.
func cacheWriteContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(ctx), 500*time.Millisecond)
}
