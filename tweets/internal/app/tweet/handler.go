package tweet

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcrosignani/uala_challenge/tweets/internal/entities"

	"github.com/labstack/echo/v4"
)

type (
	Servicier interface {
		SendTweet(ctx context.Context, tweet entities.InsertTweetRequest) error
		CreateTweet(ctx context.Context, req entities.InsertTweetRequest) (entities.Tweet, error)
		GetTweets(ctx context.Context, req entities.GetTweetsRequest) (entities.TweetsTimelineResponse, error)
	}

	Handler struct {
		service Servicier
	}
)

func NewHandler(service Servicier) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) PostTweet(c echo.Context) error {
	var req entities.InsertTweetRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err := h.service.SendTweet(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) HandleNewTweet(ctx context.Context, data []byte) error {
	var req entities.InsertTweetRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		return err
	}

	_, err = h.service.CreateTweet(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetTweets(c echo.Context) error {
	var req entities.GetTweetsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tweets, err := h.service.GetTweets(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, tweets)
}
