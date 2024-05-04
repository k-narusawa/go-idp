package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	oa "github.com/k-narusawa/go-idp/authorization/adapter"
	"github.com/k-narusawa/go-idp/authorization/adapter/gateway"
	"github.com/k-narusawa/go-idp/authorization/domain/models"
	"github.com/k-narusawa/go-idp/authorization/oauth2"
	ou "github.com/k-narusawa/go-idp/authorization/usecase"
	ra "github.com/k-narusawa/go-idp/resourceserver/adapter"
	ru "github.com/k-narusawa/go-idp/resourceserver/usecase"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
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
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	gob.Register(&models.IDSession{})
	gob.Register(&webauthn.SessionData{})
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	privateKey, err := oauth2.ReadPrivatekey()
	if err != nil {
		panic(err)
	}

	oauth2 := oauth2.NewOauth2Provider(privateKey)

	wconfig := &webauthn.Config{
		RPDisplayName: "localhost",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:3000", "http://localhost:3846"},
	}

	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		fmt.Println(err)
	}

	e.Static("/static", "static")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	db := gateway.Connect()

	isr := gateway.NewIdSessionRepository()
	ur := gateway.NewUserRepository(db)

	// oauth2
	oau := ou.NewAuthorization(oauth2, isr, ur)
	otu := ou.NewTokenUsecase(oauth2)
	oiu := ou.NewIntrospectUsecase(oauth2)
	oru := ou.NewRevokeUsecase(oauth2)
	oju := ou.NewJWKUsecase()
	olu := ou.NewLogoutUsecase(oauth2)
	owu := ou.NewWebauthnUsecase(oauth2, *webAuthn)
	wlu := ou.NewWebauthnLoginUsecase(*webAuthn)
	oa.NewOauth2Handler(e, oau, otu, oiu, oju, oru, olu, owu, wlu)

	// resource server
	uu := ru.UserinfoUsecase{}
	wu := ru.NewWebauthnUsecase(*webAuthn)
	ra.NewResourceServerHandler(e, uu, wu)

	e.Logger.Fatal(e.Start(":3846"))
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
