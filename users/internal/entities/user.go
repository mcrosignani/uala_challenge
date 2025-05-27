package entities

import "time"

type (
	User struct {
		ID        int64     `json:"id"`
		Username  string    `json:"username" validate:"required"`
		Email     string    `json:"email" validate:"required"`
		CreatedAt time.Time `json:"created_at"`
	}

	GetUserRequest struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
)
