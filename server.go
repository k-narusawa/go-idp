package main

import (
	"fmt"
	"html/template"
	"idp/authorization/adapter"
	"idp/authorization/adapter/gateway"
	"idp/authorization/oauth2"
	"idp/authorization/usecase"
	"idp/client"
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
	gateway.DbInit()

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

	privateKey, err := oauth2.ReadPrivatekey()
	if err != nil {
		panic(err)
	}

	oauth2 := oauth2.NewOauth2Provider(privateKey)
	aUsecase := usecase.NewAuthorization(oauth2)
	tUsecase := usecase.NewTokenUsecase(oauth2)
	iUsecase := usecase.NewIntrospectUsecase(oauth2)
	jUsecase := usecase.NewJWKUsecase()
	adapter.NewOauth2Handler(e, aUsecase, tUsecase, iUsecase, jUsecase)

	e.GET("/", client.IndexHandler)
	e.GET("/callback", client.CallbackHandler)
	e.GET("/userinfo", client.UserInfoHandler)

	e.Logger.Fatal(e.Start(":3846"))
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
