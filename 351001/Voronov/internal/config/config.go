package config

import (
	"Voronov/pkg/postgres"
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPport string `env:"HTTPPORT"`
	Postgres postgres.Config
}

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig("./.env", &cfg); err != nil {
		return nil, errors.New("cannot read Auth Config")
	}
	return &cfg, nil
}
