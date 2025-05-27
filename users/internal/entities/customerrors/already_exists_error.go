package customerrors

import "fmt"

var AlreadyExistsErr *AlreadyExistsError

type AlreadyExistsError struct {
	UserID     int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("user_id %d already follow user_id %d already exists", e.UserID, e.FollowerID)
}
