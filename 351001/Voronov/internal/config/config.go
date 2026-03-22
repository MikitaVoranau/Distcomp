package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPport string `env:"HTTPPORT"`
}

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig("./.env", &cfg); err != nil {
		cfg.HTTPport = "24110"
		return &cfg, nil
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.HTTPport == "" {
		return errors.New("HTTP port is required")
	}
	return nil
}
