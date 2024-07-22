package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Env string `env:"GAME_ENV" envDefault:"dev"`
	Port int `env:"PORT" envDefault:"80"`
	RedisHost string `env:"GAME_REDIS_HOST" envDefault:"127.0.0.1"`
	RedisPort int `env:"GAME_REDIS_PORT" envDefault:"36379"`
	RedisPassword string `env:"GAME_REDIS_PASSWORD" envDefault:""`
	RedisTLS string `env:"GAME_REDIS_TLS_SERVER_NAME" envDefault:""`
}

func New() (*Config,error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}