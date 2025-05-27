package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/mcrosignani/uala_challenge/users/internal/entities"
	"github.com/mcrosignani/uala_challenge/users/internal/entities/customerrors"
)

const (
	getUsersQuery = `
		SELECT id, username, email, created_at
		FROM users
	`

	countUsersQuery = `
		SELECT COUNT(*)
		FROM users
	`

	insertUserQuery = `
		INSERT INTO users (username, email, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	deleteUserQuery = `
		DELETE FROM users
		WHERE id = $1
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

func (r *Repository) GetUsers(ctx context.Context, req entities.GetUserRequest, pagination entities.Pagination) (entities.PagedResponse, error) {
	query, args := buildUsersQuery(getUsersQuery, req, pagination)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return entities.PagedResponse{}, err
	}
	defer rows.Close()

	users := make([]entities.User, 0)
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			return entities.PagedResponse{}, err
		}
		users = append(users, user)
	}

	total, err := r.Count(ctx, req)
	if err != nil {
		return entities.PagedResponse{}, err
	}

	return entities.PagedResponse{
		Total:   total,
		HasMore: total > pagination.PageSize,
		Data:    users,
	}, nil
}

func (r *Repository) Count(ctx context.Context, req entities.GetUserRequest) (int64, error) {
	query, args := buildUsersQuery(countUsersQuery, req, entities.Pagination{})

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) Create(ctx context.Context, user entities.User) (entities.User, error) {
	lastInsertId := int64(0)
	err := r.db.QueryRow(insertUserQuery, user.Username, user.Email, user.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		return entities.User{}, err
	}

	return entities.User{
		ID:        lastInsertId,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, deleteUserQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete user with id %d: %w", id, err)
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

func buildUsersQuery(query string, req entities.GetUserRequest, pagination entities.Pagination) (string, []interface{}) {
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

func getFilters(req entities.GetUserRequest) (string, []interface{}) {
	filters := []string{}
	args := make([]interface{}, 0)
	index := 1

	if req.ID != 0 {
		filters = append(filters, fmt.Sprintf("id = $%d", index))
		args = append(args, req.ID)
		index++
	}

	if req.Username != "" {
		filters = append(filters, fmt.Sprintf("username = $%d", index))
		args = append(args, req.Username)
		index++
	}

	if req.Email != "" {
		filters = append(filters, fmt.Sprintf("email = $%d", index))
		args = append(args, req.Email)
		index++
	}

	return strings.Join(filters, " AND "), args
}
