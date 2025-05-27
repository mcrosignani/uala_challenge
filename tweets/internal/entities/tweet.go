package entities

import "time"

type (
	Tweet struct {
		ID        int64     `json:"id"`
		UserID    int64     `json:"user_id" validate:"required"`
		Text      string    `json:"text" validate:"required,max=280"`
		CreatedAt time.Time `json:"created_at" validate:"required"`
	}

	InsertTweetRequest struct {
		UserID    int64     `json:"user_id" validate:"required"`
		Text      string    `json:"text" validate:"required,max=280"`
		CreatedAt time.Time `json:"-"`
	}

	TweetsTimelineRequest struct {
		Followers []int64
		After     *time.Time
		Before    *time.Time
		Limit     int
	}

	TweetsTimelineResponse struct {
		Tweets     []Tweet   `json:"tweets"`
		Count      int       `json:"count"`
		NextBefore time.Time `json:"next_before"`
	}

	GetTweetsRequest struct {
		UserID int64      `query:"user_id" validate:"required"`
		After  *time.Time `query:"after"`
		Before *time.Time `query:"before"`
		Limit  int        `query:"limit" validate:"required,min=1,max=100"`
	}
)
