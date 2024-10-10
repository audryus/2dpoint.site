package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/audryus/2dpoint.site/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server    `yaml:"server"`
		App       `yaml:"app"`
		Http      `yaml:"http"`
		Etcd      `yaml:"etcd"`
		Cockroach `yaml:"cockroach"`
	}

	Server struct {
		Header string `yaml:"header"`
		Addr   string `env-required:"true" yaml:"addr" env:"DPOINT_SERVER_ADDR"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}

	Http struct {
		Addr string `env-required:"true" yaml:"addr" env:"DPOINT_HTTP_ADDR"`
	}

	Etcd struct {
		Host string `env-required:"true" yaml:"host" env:"DPOINT_ETCD_HOST"`
		Port string `env-required:"true" yaml:"port" env:"DPOINT_ETCD_PORT"`
	}

	Cockroach struct {
		Url      string `env-required:"true" yaml:"url" env:"DPOINT_COCKROACH_URL"`
		Database string `env-required:"true" yaml:"database" env:"DPOINT_COCKROACH_DATABASE"`
		Port     string `env-required:"true" yaml:"port" env:"DPOINT_COCKROACH_PORT"`
	}
)

func New(l logger.Log) (Config, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	dir = strings.Replace(dir, "/app/", "./", 1)

	var cfg Config
	fmt.Printf("dir: %s\n", dir)
	if err := cleanenv.ReadConfig(dir+"/config.yaml", &cfg); err != nil {
		return cfg, err
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return cfg, err
	}

	l.Info("config loaded")

	return cfg, nil
}
