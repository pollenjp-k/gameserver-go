package config

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	Port       int    `env:"PORT" envDefault:"8080"`
	DBHost     string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort     int    `env:"DB_PORT" envDefault:"3306"`
	DBUser     string `env:"DB_USER" envDefault:"webapp"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"webapp_no_password"`
	DBName     string `env:"DB_NAME" envDefault:"webapp"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
