package entities

import "time"

type (
	GetFollowersRequest struct {
		UserID int64 `json:"user_id"`
	}

	Follower struct {
		UserID    int64     `json:"user_id" validate:"required"`
		FollowID  int64     `json:"follow_id" validate:"required"`
		CreatedAt time.Time `json:"created_at"`
	}

	Followers []Follower
)

func (fs Followers) GetFollowersID() []int64 {
	ids := make([]int64, len(fs))
	for _, f := range fs {
		ids = append(ids, f.FollowID)
	}
	return ids
}
