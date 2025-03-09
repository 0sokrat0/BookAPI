package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Logger   LoggerConfig   `yaml:"logger"`
}

type AppConfig struct {
	Name string `yaml:"name" env:"APP_NAME" env-default:"BookCRM"`
	Env  string `yaml:"env" env:"APP_ENV" env-default:"development"`
	Port int    `yaml:"port" env:"APP_PORT" env-default:"8080"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     uint16 `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-default:"sokrat"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"1234"`
	Name     string `yaml:"name" env:"POSTGRES_DB" env-default:"bookApi"`
	Schema   string `yaml:"schema" env:"POSTGRES_SCHEMA" env-default:"public"`
	SSLMode  string `yaml:"sslmode" env:"POSTGRES_SSLMODE" env-default:"disable"`
	MaxConn  int32  `yaml:"max_connections" env:"POSTGRES_MAX_CONN" env-default:"5"`
	MinConn  int32  `yaml:"min_connections" env:"POSTGRES_MIN_CONN" env-default:"1"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LOGGER_LEVEL" env-default:"development"`
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		if err := cleanenv.ReadConfig(".env", cfg); err != nil {
			log.Fatalf("❌ Ошибка загрузки конфигурации: %v", err)
		}

	})
	return cfg
}
