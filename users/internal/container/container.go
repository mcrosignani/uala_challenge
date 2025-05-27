package container

import (
	"github.com/mcrosignani/uala_challenge/users/internal/app/follower"
	"github.com/mcrosignani/uala_challenge/users/internal/app/user"
	"github.com/mcrosignani/uala_challenge/users/internal/config"
	"github.com/mcrosignani/uala_challenge/users/pkg/db/postgres"
)

type Dependencies struct {
	Configs config.Config

	UserHandler     *user.Handler
	FollowerHandler *follower.Handler
}

func Build() (*Dependencies, error) {
	deps := &Dependencies{}
	deps.Configs = config.NewConfig()

	db, err := postgres.NewPostgresDB(deps.Configs)
	if err != nil {
		return nil, err
	}

	if errMig := postgres.RunDBMigrations(db, deps.Configs.PostgresDB.Migrations.FilesPath); errMig != nil {
		return nil, errMig
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	deps.UserHandler = user.NewHandler(userService)

	followerRepository := follower.NewRepository(db)
	followerService := follower.NewService(followerRepository)
	deps.FollowerHandler = follower.NewHandler(followerService)

	return deps, nil
}
