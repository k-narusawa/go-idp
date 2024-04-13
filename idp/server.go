package main

import (
	"fmt"
	"html/template"
	"idp/client"
	"idp/infrastructure"
	"idp/oauth2"
	"io"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	profile, ok := os.LookupEnv("PROFILE")
	if !ok {
		fmt.Println("env is not set")
	}

	e := echo.New()
	infrastructure.DbInit()
	if profile == "local" {
		skipperFunc := func(c echo.Context) bool {
			path := c.Request().URL.Path
			return strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg")
		}
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format:  "method=${method}, uri=${uri}, status=${status}\n",
			Skipper: skipperFunc,
		}))
	} else {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())

	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.GET("/oauth2/auth", oauth2.AuthorizationEndpoint)
	e.POST("/oauth2/auth", oauth2.AuthorizationEndpoint)
	e.POST("/oauth2/token", oauth2.TokenEndpoint)
	e.POST("/oauth2/introspect", oauth2.IntrospectionEndpoint)

	e.GET("/callback", client.CallbackHandler)

	e.Logger.Fatal(e.Start(":3846"))
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
