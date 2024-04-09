package main

import (
	"html/template"
	"idp/client"
	"idp/infrastructure"
	"idp/oauth2"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	infrastructure.DbInit()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("resources/*.html")),
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
