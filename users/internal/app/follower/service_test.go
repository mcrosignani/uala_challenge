package follower_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mcrosignani/uala_challenge/users/internal/app/follower"
	"github.com/mcrosignani/uala_challenge/users/internal/entities"
	"github.com/mcrosignani/uala_challenge/users/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetFollowers(t *testing.T) {
	mockRepo := new(mocks.FollowerRepositoryMock)
	svc := follower.NewService(mockRepo)
	ctx := context.Background()
	req := entities.GetFollowersRequest{}
	pagination := entities.Pagination{}
	expected := entities.PagedResponse{Total: 2}

	mockRepo.On("GetFollowers", ctx, req, pagination).Return(expected, nil)

	resp, err := svc.GetFollowers(ctx, req, pagination)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
	mockRepo.AssertExpectations(t)
}

func TestService_InsertFollower(t *testing.T) {
	mockRepo := new(mocks.FollowerRepositoryMock)
	svc := follower.NewService(mockRepo)
	ctx := context.Background()
	req := entities.InsertFollowerRequest{UserID: 1, FollowID: 2}

	// Usamos mock.MatchedBy para ignorar el valor exacto de CreatedAt
	mockRepo.On("InsertFollower", ctx, mock.MatchedBy(func(r entities.InsertFollowerRequest) bool {
		return r.UserID == req.UserID && r.FollowID == req.FollowID && !r.CreatedAt.IsZero()
	})).Return(nil)

	err := svc.InsertFollower(ctx, req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_InsertFollower_Error(t *testing.T) {
	mockRepo := new(mocks.FollowerRepositoryMock)
	svc := follower.NewService(mockRepo)
	ctx := context.Background()
	req := entities.InsertFollowerRequest{UserID: 1, FollowID: 2}

	mockRepo.On("InsertFollower", ctx, mock.Anything).Return(errors.New("insert error"))

	err := svc.InsertFollower(ctx, req)
	assert.Error(t, err)
}

func TestService_DeleteFollower(t *testing.T) {
	mockRepo := new(mocks.FollowerRepositoryMock)
	svc := follower.NewService(mockRepo)
	ctx := context.Background()
	userID := int64(1)
	followID := int64(2)

	mockRepo.On("DeleteFollower", ctx, entities.DeleteFollowerRequest{UserID: userID, FollowID: followID}).Return(nil)

	err := svc.DeleteFollower(ctx, userID, followID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteFollower_Error(t *testing.T) {
	mockRepo := new(mocks.FollowerRepositoryMock)
	svc := follower.NewService(mockRepo)
	ctx := context.Background()
	userID := int64(1)
	followID := int64(2)

	mockRepo.On("DeleteFollower", ctx, entities.DeleteFollowerRequest{UserID: userID, FollowID: followID}).Return(errors.New("delete error"))

	err := svc.DeleteFollower(ctx, userID, followID)
	assert.Error(t, err)
}
