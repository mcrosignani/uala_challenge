package tweet

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/mcrosignani/uala_challenge/tweets/internal/entities"
	"github.com/mcrosignani/uala_challenge/tweets/pkg/clients"
	"github.com/mcrosignani/uala_challenge/tweets/pkg/nats"
)

type (
	Repositorier interface {
		CreateTweet(ctx context.Context, req entities.InsertTweetRequest) (entities.Tweet, error)
		TweetsTimeline(ctx context.Context, req entities.TweetsTimelineRequest) (entities.TweetsTimelineResponse, error)
	}

	Service struct {
		messajeService nats.MessageService
		repository     Repositorier
		usersClient    clients.UsersClient
	}
)

func NewService(messageService nats.MessageService, repository Repositorier, usersClient clients.UsersClient) *Service {
	return &Service{
		messajeService: messageService,
		repository:     repository,
		usersClient:    usersClient,
	}
}

func (s *Service) SendTweet(ctx context.Context, tweet entities.InsertTweetRequest) error {
	tweetMsg, err := json.Marshal(tweet)
	if err != nil {
		return err
	}

	return s.messajeService.Publish("tweets", tweetMsg)
}

func (s *Service) CreateTweet(ctx context.Context, req entities.InsertTweetRequest) (entities.Tweet, error) {
	req.CreatedAt = time.Now().UTC()
	return s.repository.CreateTweet(ctx, req)
}

func (s *Service) GetTweets(ctx context.Context, req entities.GetTweetsRequest) (entities.TweetsTimelineResponse, error) {
	followersResp, err := s.usersClient.GetFollowersByUserID(ctx, entities.GetFollowersRequest{
		UserID: req.UserID,
	})
	if err != nil {
		return entities.TweetsTimelineResponse{}, err
	}

	var followers entities.Followers
	bytes, _ := json.Marshal(followersResp.Data)
	err = json.Unmarshal(bytes, &followers)
	if err != nil {
		return entities.TweetsTimelineResponse{}, errors.New("cannot parse followers data")
	}

	return s.repository.TweetsTimeline(ctx, entities.TweetsTimelineRequest{
		Followers: followers.GetFollowersID(),
		After:     req.After,
		Before:    req.Before,
		Limit:     req.Limit,
	})
}
