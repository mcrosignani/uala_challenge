package container

import (
	"github.com/mcrosignani/uala_challenge/tweets/internal/app/tweet"
	"github.com/mcrosignani/uala_challenge/tweets/internal/config"
	"github.com/mcrosignani/uala_challenge/tweets/pkg/clients"
	"github.com/mcrosignani/uala_challenge/tweets/pkg/db/postgres"
	"github.com/mcrosignani/uala_challenge/tweets/pkg/nats"
)

type Dependencies struct {
	Configs        config.Config
	MessageService nats.MessageService

	TweetHandler *tweet.Handler
}

func Build() (*Dependencies, error) {
	deps := &Dependencies{}
	deps.Configs = config.NewConfig()

	messajeService, err := nats.NewNCService(deps.Configs)
	if err != nil {
		return nil, err
	}
	deps.MessageService = messajeService

	db, err := postgres.NewPostgresDB(deps.Configs)
	if err != nil {
		return nil, err
	}

	if errMig := postgres.RunDBMigrations(db, deps.Configs.PostgresDB.Migrations.FilesPath); errMig != nil {
		return nil, errMig
	}

	usersClient := clients.NewUsersClient(deps.Configs)

	tweetRepository := tweet.NewRepository(db)
	tweetService := tweet.NewService(messajeService, tweetRepository, usersClient)
	deps.TweetHandler = tweet.NewHandler(tweetService)

	return deps, nil
}
