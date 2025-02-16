package config

import (
	// "time"

	// "github.com/cenkalti/backoff/v4"
)

type Config struct {
	Workers         int           `mapstructure:"workers"`
	MaxQueueSize    int           `mapstructure:"max_queue_size"`
	Redis           RedisConfig   `mapstructure:"redis"`
	MetricsPort     int           `mapstructure:"metrics_port"`
	RateLimit       int           `mapstructure:"rate_limit"`
	DefaultRetries  int           `mapstructure:"default_retries"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func Load() (*Config, error) {
	// Implement configuration loading from file/env
	// Using default values for demonstration
	return &Config{
		Workers:        100,
		MaxQueueSize:  1e6,
		MetricsPort:   9090,
		RateLimit:     1000,
		DefaultRetries: 3,
		Redis: RedisConfig{
			Address: "localhost:6379",
		},
	}, nil
}