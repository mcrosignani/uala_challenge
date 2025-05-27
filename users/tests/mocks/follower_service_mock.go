package mocks

import (
	"context"

	"github.com/mcrosignani/uala_challenge/users/internal/entities"
	"github.com/stretchr/testify/mock"
)

type FollowerServiceMock struct {
	mock.Mock
}

func (m *FollowerServiceMock) GetFollowers(ctx context.Context, req entities.GetFollowersRequest, pagination entities.Pagination) (entities.PagedResponse, error) {
	args := m.Called(ctx, req, pagination)
	return args.Get(0).(entities.PagedResponse), args.Error(1)
}
func (m *FollowerServiceMock) InsertFollower(ctx context.Context, req entities.InsertFollowerRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
func (m *FollowerServiceMock) DeleteFollower(ctx context.Context, userID, followID int64) error {
	args := m.Called(ctx, userID, followID)
	return args.Error(0)
}
