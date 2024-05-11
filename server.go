package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"

	oa "github.com/k-narusawa/go-idp/authorization/adapter"
	"github.com/k-narusawa/go-idp/authorization/adapter/gateway"
	"github.com/k-narusawa/go-idp/authorization/domain/models"
	"github.com/k-narusawa/go-idp/authorization/oauth2"
	ou "github.com/k-narusawa/go-idp/authorization/usecase"
	"github.com/k-narusawa/go-idp/logger"
	ra "github.com/k-narusawa/go-idp/resources/adapter"
	ru "github.com/k-narusawa/go-idp/resources/usecase"

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

	_, ok := os.LookupEnv("PROFILE")
	if !ok {
		fmt.Println("env is not set")
	}

	e := echo.New()
	gateway.DbInit()

	logger := logger.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Info("REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.Warn("REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	gob.Register(&models.IdpSession{})
	gob.Register(&webauthn.SessionData{})
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Static("/static", "static")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	db := gateway.Connect()

	privateKey, err := oauth2.ReadPrivatekey()
	if err != nil {
		panic(err)
	}

	wconfig := &webauthn.Config{
		RPDisplayName: "localhost",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:3000", "http://localhost:3846"},
	}

	oauth2 := oauth2.NewOauth2Provider(privateKey, *logger)
	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		panic(err)
	}

	ur := gateway.NewUserRepository(db)
	wcr := gateway.NewWebauthnCredentialRepository(db)
	wsr := gateway.NewWebauthnSessionRepository(db)
	cr := gateway.NewClientRepository(db)
	isr := gateway.NewIdpSessionRepository()
	osr := gateway.NewOidcSessionRepository(db)
	lssr := gateway.NewLoginSkipSessionRepository(db)
	atr := gateway.NewAccessTokenRepository(db)

	// oauth2
	oau := ou.NewAuthorization(oauth2, ur, isr, osr, lssr)
	otu := ou.NewTokenUsecase(oauth2, atr)
	oiu := ou.NewIntrospectUsecase(oauth2)
	oru := ou.NewRevokeUsecase(oauth2)
	oju := ou.NewJWKUsecase()
	olu := ou.NewLogoutUsecase(oauth2, isr, osr)
	osu := ou.NewSessionUsecase(ur, isr, lssr)
	owu := ou.NewAuthenticateWebauthnUsecase(oauth2, *webAuthn, ur, wcr, lssr)
	oa.NewOauth2Handler(e, oau, otu, oiu, oju, oru, olu, osu, owu)

	// client
	cu := ou.NewClientUsecase(cr)
	oa.NewClientHandler(e, cu)

	// resource server
	uu := ru.UserinfoUsecase{}
	wu := ru.NewWebauthnUsecase(*logger, *webAuthn, ur, wcr, wsr)
	iu := ru.NewIntrospectUsecase(*logger, atr)
	ra.NewResourceServerHandler(e, uu, wu, iu)

	e.Logger.Fatal(e.Start(":3846"))
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
