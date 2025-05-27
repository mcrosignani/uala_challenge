package mocks

import (
	"context"

	"github.com/mcrosignani/uala_challenge/users/internal/entities"
	"github.com/stretchr/testify/mock"
)

type FollowerRepositoryMock struct {
	mock.Mock
}

func (m *FollowerRepositoryMock) GetFollowers(ctx context.Context, req entities.GetFollowersRequest, pagination entities.Pagination) (entities.PagedResponse, error) {
	args := m.Called(ctx, req, pagination)
	return args.Get(0).(entities.PagedResponse), args.Error(1)
}
func (m *FollowerRepositoryMock) InsertFollower(ctx context.Context, req entities.InsertFollowerRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
func (m *FollowerRepositoryMock) DeleteFollower(ctx context.Context, req entities.DeleteFollowerRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
