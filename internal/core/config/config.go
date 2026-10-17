package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	TimeZone *time.Location
}

const defaultTZ = "UTC"

func NewConfig() (*Config, error) {
	timeZone := os.Getenv("TIME_ZONE")
	if timeZone == "" {
		timeZone = defaultTZ
	}

	zone, err := time.LoadLocation(timeZone)
	if err != nil {
		return nil, fmt.Errorf("load time zone '%s': %w", timeZone, err)
	}

	return &Config{
		TimeZone: zone,
	}, nil
}

func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get core config: %w", err))
	}

	return config
}
