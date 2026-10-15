package tasks_postgres_repository

import core_postgres_pool "github.com/darya-kznv/todo-service/internal/core/repository/postgres/pool"

type TasksRepository struct {
	pool core_postgres_pool.Pool
}

func NewTasksRepository(p core_postgres_pool.Pool) *TasksRepository {
	return &TasksRepository{
		pool: p,
	}
}
