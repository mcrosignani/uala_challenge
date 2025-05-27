package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/mcrosignani/uala_challenge/tweets/internal/config"
	"github.com/mcrosignani/uala_challenge/tweets/internal/entities"
)

type (
	UsersClient interface {
		GetFollowersByUserID(ctx context.Context, req entities.GetFollowersRequest) (entities.PagedResponse, error)
	}

	userClient struct {
		cfg        config.Config
		httpClient *resty.Client
	}
)

func NewUsersClient(cfg config.Config) UsersClient {
	return &userClient{
		cfg: cfg,
		httpClient: resty.New().
			SetTimeout(cfg.UsersClient.Timeout).
			SetRetryCount(cfg.UsersClient.RetryCount),
	}
}

func (c *userClient) GetFollowersByUserID(ctx context.Context, req entities.GetFollowersRequest) (entities.PagedResponse, error) {
	url := fmt.Sprintf("http://%s:%s/followers?user_id=%d", c.cfg.UsersClient.Host, c.cfg.UsersClient.Port, req.UserID)
	resp, err := c.httpClient.R().
		SetContext(ctx).
		Get(url)
	if err != nil {
		return entities.PagedResponse{}, err
	}

	if resp.IsError() {
		return entities.PagedResponse{}, fmt.Errorf("error response: status %d, body: %s", resp.StatusCode(), resp.Body())
	}

	var followers entities.PagedResponse
	if err := json.Unmarshal(resp.Body(), &followers); err != nil {
		return entities.PagedResponse{}, err
	}

	return followers, nil
}
