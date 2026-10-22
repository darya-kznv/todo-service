package tasks_cached_repository

import (
	core_logger "github.com/darya-kznv/todo-service/internal/core/logger"
	core_redis_pool "github.com/darya-kznv/todo-service/internal/core/repository/redis/pool"
	tasks_service "github.com/darya-kznv/todo-service/internal/features/tasks/service"
)

type CachedRepository struct {
	pool           core_redis_pool.Pool
	mainRepository tasks_service.TasksRepository
	logger         *core_logger.Logger
}

func NewCachedRepository(
	pool core_redis_pool.Pool,
	mainRepository tasks_service.TasksRepository,
	logger *core_logger.Logger,
) *CachedRepository {
	return &CachedRepository{
		pool:           pool,
		mainRepository: mainRepository,
		logger:         logger,
	}
}
