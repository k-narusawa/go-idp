package main

import (
	"github.com/GIT_USER_ID/GIT_REPO_ID/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	//todo: handle the error!
	c, _ := handlers.NewContainer()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())


	// WebauthnStart - 登録開始
	e.GET("/api/v1/webauthn/start", c.WebauthnStart)


	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
