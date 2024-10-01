package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP `yaml:"http"`
	}

	HTTP struct {
		Addr string `env-required:"true" yaml:"addr" env:"HTTP_ADDR"`
	}
)

func New() (Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("./config.yaml", &cfg); err != nil {
		return cfg, err
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
