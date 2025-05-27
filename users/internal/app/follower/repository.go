package follower

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/mcrosignani/uala_challenge/users/internal/entities"
	"github.com/mcrosignani/uala_challenge/users/internal/entities/customerrors"
)

const (
	selectFollowersQuery = `
		SELECT user_id, follow_id, created_at
		FROM followers
		`

	countFollowersQuery = `
		SELECT COUNT(*)
		FROM followers
	`

	insertFollowerQuery = `
		INSERT INTO followers (user_id, follow_id, created_at) 
		VALUES ($1, $2, $3)
	`

	deleteFollowerQuery = `
		DELETE FROM followers
		WHERE user_id = $1 AND follow_id = $2
	`
)

type (
	Repository struct {
		db *sql.DB
	}
)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetFollowers(ctx context.Context, req entities.GetFollowersRequest, pagination entities.Pagination) (entities.PagedResponse, error) {
	query, args := buildFollowerQuery(selectFollowersQuery, req, pagination)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return entities.PagedResponse{}, err
	}
	defer rows.Close()

	follows := make([]entities.Follower, 0)
	for rows.Next() {
		var follower entities.Follower
		if err := rows.Scan(&follower.UserID, &follower.FollowID, &follower.CreatedAt); err != nil {
			return entities.PagedResponse{}, err
		}
		follows = append(follows, follower)
	}

	total, err := r.Count(ctx, req)
	if err != nil {
		return entities.PagedResponse{}, err
	}

	return entities.PagedResponse{
		Total:   total,
		HasMore: total > pagination.PageSize,
		Data:    follows,
	}, nil
}

func (r *Repository) Count(ctx context.Context, req entities.GetFollowersRequest) (int64, error) {
	query, args := buildFollowerQuery(countFollowersQuery, req, entities.Pagination{})

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) InsertFollower(ctx context.Context, req entities.InsertFollowerRequest) error {
	_, err := r.db.ExecContext(ctx, insertFollowerQuery, req.UserID, req.FollowID, req.CreatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return &customerrors.AlreadyExistsError{
					UserID:     req.UserID,
					FollowerID: req.FollowID,
				}
			case "23503":
				return &customerrors.NotFoundError{}
			}
		}
		return err
	}

	return nil
}

func (r *Repository) DeleteFollower(ctx context.Context, req entities.DeleteFollowerRequest) error {
	result, err := r.db.ExecContext(ctx, deleteFollowerQuery, req.UserID, req.FollowID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return &customerrors.NotFoundError{}
			}
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &customerrors.NotFoundError{}
	}

	return nil
}

func buildFollowerQuery(query string, req entities.GetFollowersRequest, pagination entities.Pagination) (string, []interface{}) {
	filters, args := getFilters(req)
	if filters != "" {
		query += fmt.Sprintf(" WHERE %s", filters)
	}

	if pagination.PageSize > 0 {
		offset := pagination.PageSize * (pagination.SafePageNumber() - 1)
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.PageSize, offset)
	}

	return query, args
}

func getFilters(req entities.GetFollowersRequest) (string, []interface{}) {
	filters := []string{}
	args := make([]interface{}, 0)
	index := 1

	if req.UserID != 0 {
		filters = append(filters, fmt.Sprintf("user_id = $%d", index))
		args = append(args, req.UserID)
		index++
	}

	return strings.Join(filters, " AND "), args
}
