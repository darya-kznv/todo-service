package core_goredis_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" default:"6379"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	DB       int           `envconfig:"DB" default:"0"`
	TTL      time.Duration `envconfig:"TTL" default:"5m"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("REDIS", &config); err != nil {
		return Config{}, fmt.Errorf("process redis envconfig: %w", err)
	}

	if config.TTL <= 0 {
		return Config{}, fmt.Errorf("redis TTL must be positive")
	}

	if config.DB < 0 {
		return Config{}, fmt.Errorf("redis DB must be non-negative")
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()

	if err != nil {
		panic(fmt.Errorf("get Redis connection pool config: %w", err))
	}

	return config
}
