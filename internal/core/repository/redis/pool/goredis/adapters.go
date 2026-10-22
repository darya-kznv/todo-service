package core_goredis_pool

import (
	"errors"

	core_redis_pool "github.com/darya-kznv/todo-service/internal/core/repository/redis/pool"
	"github.com/redis/go-redis/v9"
)

func mapError(err error) error {
	if errors.Is(err, redis.Nil) {
		return core_redis_pool.NotFound
	}

	return err
}

type stringCmd struct{ command *redis.StringCmd }

func (c stringCmd) Result() (string, error) {
	value, err := c.command.Result()

	return value, mapError(err)
}

func (c stringCmd) Err() error {
	return mapError(c.command.Err())
}

type statusCmd struct{ command *redis.StatusCmd }

func (c statusCmd) Result() (string, error) {
	value, err := c.command.Result()

	return value, mapError(err)
}

func (c statusCmd) Err() error {
	return mapError(c.command.Err())
}

type intCmd struct{ command *redis.IntCmd }

func (c intCmd) Result() (int64, error) {
	value, err := c.command.Result()
	return value, mapError(err)
}

func (c intCmd) Err() error {
	return mapError(c.command.Err())
}
