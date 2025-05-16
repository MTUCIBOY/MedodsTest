package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`

	ServerSetting `yaml:"http_server"`
}

type ServerSetting struct {
	Address      string        `yaml:"address"       env-default:"localhost:8888"`
	ReadTimeout  time.Duration `yaml:"read_timeout"  env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"5s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"  env-default:"10s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed parse config file: " + err.Error())
	}

	return &cfg
}
