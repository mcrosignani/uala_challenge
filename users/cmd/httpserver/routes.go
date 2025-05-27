package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/mcrosignani/uala_challenge/users/internal/container"
)

func SetRoutes(root *echo.Group, deps *container.Dependencies) {
	usersGroup := root.Group("users")
	usersGroup.GET("", deps.UserHandler.GetUsers)
	usersGroup.POST("", deps.UserHandler.Create)
	usersGroup.DELETE("/:id", deps.UserHandler.Delete)

	followerGroup := root.Group("followers")
	followerGroup.GET("", deps.FollowerHandler.GetFollowers)
	followerGroup.POST("/:user_id", deps.FollowerHandler.Create)
	followerGroup.DELETE("/:user_id/delete/:follow_id", deps.FollowerHandler.Delete)
}
