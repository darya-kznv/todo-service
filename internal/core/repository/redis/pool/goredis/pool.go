package core_goredis_pool

import (
	"context"
	"fmt"
	"net"
	"time"

	core_redis_pool "github.com/darya-kznv/todo-service/internal/core/repository/redis/pool"
	"github.com/redis/go-redis/v9"
)

type Pool struct {
	client *redis.Client
	ttl    time.Duration
}

var _ core_redis_pool.Pool = (*Pool)(nil)

func NewPool(ctx context.Context, config Config) (*Pool, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Pool{client: client, ttl: config.TTL}, nil
}

func (p *Pool) Get(ctx context.Context, key string) core_redis_pool.StringCmd {
	return stringCmd{p.client.Get(ctx, key)}
}

func (p *Pool) Set(ctx context.Context, key string, value any, ttl time.Duration) core_redis_pool.StatusCmd {
	return statusCmd{p.client.Set(ctx, key, value, ttl)}
}

func (p *Pool) Del(ctx context.Context, keys ...string) core_redis_pool.IntCmd {
	return intCmd{p.client.Del(ctx, keys...)}
}

func (p *Pool) HGet(ctx context.Context, key string, field string) core_redis_pool.StringCmd {
	return stringCmd{p.client.HGet(ctx, key, field)}
}

func (p *Pool) HSet(ctx context.Context, key string, values ...any) core_redis_pool.IntCmd {
	return intCmd{p.client.HSet(ctx, key, values...)}
}

func (p *Pool) Expire(ctx context.Context, key string, ttl time.Duration) core_redis_pool.BoolCmd {
	return boolCmd{p.client.Expire(ctx, key, ttl)}
}

func (p *Pool) Close() error { return p.client.Close() }

func (p *Pool) TTL() time.Duration { return p.ttl }
