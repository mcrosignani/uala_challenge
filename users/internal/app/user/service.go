package user

import (
	"context"
	"time"

	"github.com/mcrosignani/uala_challenge/users/internal/entities"
)

type (
	Repositorier interface {
		GetUsers(ctx context.Context, req entities.GetUserRequest, pagination entities.Pagination) (entities.PagedResponse, error)
		Create(ctx context.Context, user entities.User) (entities.User, error)
		Delete(ctx context.Context, id int64) error
	}

	Serivce struct {
		repository Repositorier
	}
)

func NewService(repository Repositorier) *Serivce {
	return &Serivce{
		repository: repository,
	}
}

func (s *Serivce) GetUsers(ctx context.Context, req entities.GetUserRequest, pagination entities.Pagination) (entities.PagedResponse, error) {
	return s.repository.GetUsers(ctx, req, pagination)
}

func (s *Serivce) Create(ctx context.Context, user entities.User) (entities.User, error) {
	user.CreatedAt = time.Now().UTC()
	return s.repository.Create(ctx, user)
}

func (s *Serivce) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
