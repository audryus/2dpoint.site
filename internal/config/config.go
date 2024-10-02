package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server `yaml:"server"`
		App    `yaml:"app"`
		HTTP   `yaml:"http"`
	}

	Server struct {
		Header string `yaml:"header"`
		Addr   string `env-required:"true" yaml:"addr" env:"SERVER_ADDR"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
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
