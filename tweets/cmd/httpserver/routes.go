package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/mcrosignani/uala_challenge/tweets/internal/container"
)

func SetRoutes(root *echo.Group, deps *container.Dependencies) {
	tweetGroup := root.Group("/tweets")
	tweetGroup.POST("", deps.TweetHandler.PostTweet)
	tweetGroup.GET("/timeline", deps.TweetHandler.GetTweets)
}
