package tweet

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/mcrosignani/uala_challenge/tweets/internal/entities"
)

const (
	selectTweetsQuery = `
		SELECT id, user_id, text, created_at
		FROM tweets
	`

	insertTweetQuery = `
		INSERT INTO tweets (user_id, text, created_at) 
		VALUES ($1, $2, $3)
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

func (r *Repository) TweetsTimeline(ctx context.Context, req entities.TweetsTimelineRequest) (entities.TweetsTimelineResponse, error) {
	query, args := buildTweetsQuery(selectTweetsQuery, req)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return entities.TweetsTimelineResponse{}, err
	}
	defer rows.Close()

	tweets := make([]entities.Tweet, 0)
	for rows.Next() {
		var tweet entities.Tweet
		if err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Text, &tweet.CreatedAt); err != nil {
			return entities.TweetsTimelineResponse{}, err
		}
		tweets = append(tweets, tweet)
	}

	lastTweet := tweets[len(tweets)-1]
	return entities.TweetsTimelineResponse{
		Tweets:     tweets,
		Count:      len(tweets),
		NextBefore: lastTweet.CreatedAt,
	}, nil
}

func (r *Repository) CreateTweet(ctx context.Context, req entities.InsertTweetRequest) (entities.Tweet, error) {
	lastInsertId := int64(0)
	err := r.db.QueryRow(insertTweetQuery, req.UserID, req.Text, req.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		return entities.Tweet{}, err
	}

	return entities.Tweet{
		ID:        lastInsertId,
		UserID:    req.UserID,
		Text:      req.Text,
		CreatedAt: req.CreatedAt,
	}, nil
}

func buildTweetsQuery(query string, req entities.TweetsTimelineRequest) (string, []interface{}) {
	filters, args := getFilters(req)
	if filters != "" {
		query += fmt.Sprintf(" WHERE %s", filters)
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT %d", req.Limit)

	return query, args
}

func getFilters(req entities.TweetsTimelineRequest) (string, []interface{}) {
	filters := []string{}
	args := make([]interface{}, 0)
	index := 1

	if len(req.Followers) > 0 {
		placeholders := make([]string, len(req.Followers))
		for i := range req.Followers {
			placeholders[i] = fmt.Sprintf("$%d", index)
			args = append(args, req.Followers[i])
			index++
		}
		filters = append(filters, fmt.Sprintf("user_id IN (%s)", strings.Join(placeholders, ", ")))
	}

	if req.After != nil {
		filters = append(filters, fmt.Sprintf("created_at > $%d", index))
		args = append(args, *req.Before)
		index++
	}

	if req.Before != nil {
		filters = append(filters, fmt.Sprintf("created_at < $%d", index))
		args = append(args, *req.Before)
		index++
	}

	return strings.Join(filters, " AND "), args
}
