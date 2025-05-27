package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mcrosignani/uala_challenge/users/cmd/httpserver"
	"github.com/mcrosignani/uala_challenge/users/internal/container"
)

func main() {
	deps, err := container.Build()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Validator = httpserver.NewValidator()

	httpserver.SetRoutes(e.Group("/"), deps)
	e.Logger.Fatal(e.Start(":" + deps.Configs.Port))
}
