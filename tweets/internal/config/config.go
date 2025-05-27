package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ProjectName string `envconfig:"PROJECT_NAME" default:"tweets"`
	Port        string `envconfig:"PORT" default:"8090"`
	Nats        struct {
		Host string `envconfig:"NATS_HOST" default:"localhost"`
		Port int    `envconfig:"NATS_PORT" default:"4222"`
	}
	PostgresDB struct {
		Host       string `envconfig:"POSTGRES_HOST" default:"localhost"`
		Port       string `envconfig:"POSTGRES_PORT" default:"5532"`
		User       string `envconfig:"POSTGRES_USER" default:"postgres"`
		Password   string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
		DBName     string `envconfig:"POSTGRES_DB" default:"tweets_db"`
		Migrations struct {
			FilesPath string `envconfig:"POSTGRES_MIGRATIONS_PATH" default:"../../../db/migrations"`
		}
	}
	UsersClient struct {
		Host       string        `envconfig:"USERS_HOST" default:"localhost"`
		Port       string        `envconfig:"USERS_PORT" default:"8080"`
		Timeout    time.Duration `envconfig:"USERS_TIMEOUT" default:"1s"`
		RetryCount int           `envconfig:"USERS_RETRY_COUNT" default:"3"`
	}
}

func NewConfig() Config {
	configs := Config{}
	err := envconfig.Process("", &configs)
	if err != nil {
		panic(err)
	}

	return configs
}
