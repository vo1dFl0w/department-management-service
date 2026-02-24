package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	Username string `yaml:"username" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `yaml:"dbname" env:"POSTGRES_DB" env-required:"true"`
	Sslmode  string `yaml:"sslmode" env:"POSTGRES_SSL" env-required:"true"`
}

type Config struct {
	Env              string         `yaml:"env" env:"SERVER_ENV" env-required:"true"`
	HTTPAddr         string         `yaml:"http_addr" env:"SERVER_HTTP_ADDR" env-required:"true"`
	RequestTimeout   time.Duration  `yaml:"request_timeout"`
	LoggerTimeFormat string         `yaml:"logger_time_format"`
	ShutdownTimeout  time.Duration  `yaml:"shutdown_timeout"`
	DB               DatabaseConfig `yaml:"db"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return nil, fmt.Errorf("CONFIG_PATH not set")
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	return &cfg, nil
}
