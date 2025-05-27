package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ProjectName string `envconfig:"PROJECT_NAME" default:"users"`
	Port        string `envconfig:"PORT" default:"8080"`
	PostgresDB  struct {
		Host       string `envconfig:"POSTGRES_HOST" default:"localhost"`
		Port       string `envconfig:"POSTGRES_PORT" default:"5432"`
		User       string `envconfig:"POSTGRES_USER" default:"postgres"`
		Password   string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
		DBName     string `envconfig:"POSTGRES_DB" default:"users_db"`
		Migrations struct {
			FilesPath string `envconfig:"POSTGRES_MIGRATIONS_PATH" default:"../../../db/migrations"`
		}
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
