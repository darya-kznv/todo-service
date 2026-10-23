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

type stringCmd struct {
	command *redis.StringCmd
}

func (c stringCmd) Result() (string, error) {
	value, err := c.command.Result()

	return value, mapError(err)
}

func (c stringCmd) Err() error {
	return mapError(c.command.Err())
}

func (c stringCmd) Bytes() ([]byte, error) {
	value, err := c.command.Bytes()

	return value, mapError(err)
}

type statusCmd struct {
	command *redis.StatusCmd
}

func (c statusCmd) Result() (string, error) {
	value, err := c.command.Result()

	return value, mapError(err)
}

func (c statusCmd) Err() error {
	return mapError(c.command.Err())
}

type intCmd struct {
	command *redis.IntCmd
}

func (c intCmd) Result() (int64, error) {
	value, err := c.command.Result()

	return value, mapError(err)
}

func (c intCmd) Err() error {
	return mapError(c.command.Err())
}

type boolCmd struct {
	command *redis.BoolCmd
}

func (c boolCmd) Result() (bool, error) {
	value, err := c.command.Result()

	return value, mapError(err)
}

func (c boolCmd) Err() error {
	return mapError(c.command.Err())
}
