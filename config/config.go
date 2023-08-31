package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
)

type Config struct {
	App struct {
		Name    string `yaml:"name" env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	} `yaml:"app"`
	Log struct {
		Level    string `yaml:"level" env:"LOG_LEVE"`
		LogsPath string `yaml:"logs_path" env:"LOG_PATH"`
	}
	HTTP struct {
		BindIP string `yaml:"bind_ip" env:"HTTP_BIND_IP"`
		Port   string `yaml:"port" env:"HTTP_PORT"`
	} `yaml:"http"`
	PostgreSQL struct {
		Host        string `yaml:"host" env:"POSTGRE_HOST"`
		Port        string `yaml:"port" env:"POSTGRE_PORT"`
		Username    string `yaml:"username" env:"POSTGRE_USERNAME"`
		Password    string `yaml:"password" env:"POSTGRE_PASSWORD"`
		Database    string `yaml:"database" env:"POSTGRE_DATABASE"`
		MaxPoolSize int    `yaml:"max_pool_size" env:"POSTGRE_MAX_POOL_SIZE"`
	} `yaml:"postgresql"`
}

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(path.Join("./", configPath), cfg); err != nil {
		return nil, fmt.Errorf("failed to read config due to error: %w", err)
	}

	if err := cleanenv.UpdateEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to update environment variables due to error: %w", err)
	}

	return cfg, nil
}
