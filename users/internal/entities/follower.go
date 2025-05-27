package entities

import "time"

type (
	Follower struct {
		UserID    int64     `json:"user_id" validate:"required"`
		FollowID  int64     `json:"follow_id" validate:"required"`
		CreatedAt time.Time `json:"created_at"`
	}

	GetFollowersRequest struct {
		UserID int64 `query:"user_id"`
	}

	InsertFollowerRequest struct {
		UserID    int64     `param:"user_id" json:"user_id" validate:"required"`
		FollowID  int64     `json:"follow_id" validate:"required"`
		CreatedAt time.Time `json:"-"`
	}

	DeleteFollowerRequest struct {
		UserID   int64 `param:"user_id" json:"user_id" validate:"required"`
		FollowID int64 `param:"follow_id" json:"follow_id" validate:"required"`
	}
)
