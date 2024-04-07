package main

import (
	"idp/oauth2"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/oauth2/auth", oauth2.AuthorizationEndpoint)
	e.POST("/oauth2/auth", oauth2.AuthorizationEndpoint)
	e.POST("/oauth2/token", oauth2.TokenEndpoint)
	// e.HandleFunc("/oauth2/introspect", oauth2.IntrospectionEndpoint)

	e.Logger.Fatal(e.Start(":3846"))
}
