package follower

import (
	"context"
	"time"

	"github.com/mcrosignani/uala_challenge/users/internal/entities"
)

type (
	Repositorier interface {
		GetFollowers(ctx context.Context, req entities.GetFollowersRequest, pagination entities.Pagination) (entities.PagedResponse, error)
		InsertFollower(ctx context.Context, req entities.InsertFollowerRequest) error
		DeleteFollower(ctx context.Context, req entities.DeleteFollowerRequest) error
	}

	Service struct {
		repository Repositorier
	}
)

func NewService(repository Repositorier) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetFollowers(ctx context.Context, req entities.GetFollowersRequest, pagination entities.Pagination) (entities.PagedResponse, error) {
	return s.repository.GetFollowers(ctx, req, pagination)
}

func (s *Service) InsertFollower(ctx context.Context, req entities.InsertFollowerRequest) error {
	req.CreatedAt = time.Now().UTC()
	return s.repository.InsertFollower(ctx, req)
}

func (s *Service) DeleteFollower(ctx context.Context, userID, followID int64) error {
	return s.repository.DeleteFollower(ctx, entities.DeleteFollowerRequest{
		UserID:   userID,
		FollowID: followID,
	})
}
