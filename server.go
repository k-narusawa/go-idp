package main

import (
	"encoding/gob"
	"html/template"
	"io"

	"os"

	"github.com/k-narusawa/go-idp/adapter/controllers"
	"github.com/k-narusawa/go-idp/adapter/gateways"
	oa "github.com/k-narusawa/go-idp/authorization/application"
	"github.com/k-narusawa/go-idp/authorization/oauth2"
	"github.com/k-narusawa/go-idp/cert"
	"github.com/k-narusawa/go-idp/domain/models"
	"github.com/k-narusawa/go-idp/logger"
	gmiddleware "github.com/k-narusawa/go-idp/middleware"
	ra "github.com/k-narusawa/go-idp/resources/application"
	"gopkg.in/yaml.v2"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Mode   string `yaml:"mode"`
	Server struct {
		Port    string `yaml:"port"`
		Session struct {
			Secret string `yaml:"secret"`
		} `yaml:"session"`
		Cors struct {
			AllowedOrigins   []string `yaml:"allowedOrigins"`
			AllowMethods     []string `yaml:"allowMethods"`
			AllowHeaders     []string `yaml:"allowHeaders"`
			AllowCredentials bool     `yaml:"allowCredentials"`
			ExposeHeaders    []string `yaml:"exposeHeaders"`
			MaxAge           int      `yaml:"maxAge"`
		} `yaml:"cors"`
	} `yaml:"server"`
	DB struct {
		Mode string `yaml:"mode"`
		DSN  string `yaml:"dsn"`
	} `yaml:"db"`
	Webauthn struct {
		DisplayName string   `yaml:"displayName"`
		RPID        string   `yaml:"rpId"`
		RPOrigins   []string `yaml:"rpOrigins"`
	} `yaml:"webauthn"`
}

type TemplateRenderer struct {
	templates *template.Template
}

func main() {
	config := loadConfig()

	e := echo.New()
	gateways.DbInit(
		/* mode     */ config.DB.Mode,
		/* dsn      */ config.DB.DSN,
	)

	logger := logger.New()

	e.Use(gmiddleware.NewLogger(*logger))
	e.Use(middleware.Recover())

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.Server.Session.Secret))))
	gob.Register(&models.IdpSession{})
	gob.Register(&webauthn.SessionData{})

	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Static("/static", "static")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.Server.Cors.AllowedOrigins,
		AllowMethods:     config.Server.Cors.AllowMethods,
		AllowHeaders:     config.Server.Cors.AllowHeaders,
		AllowCredentials: config.Server.Cors.AllowCredentials,
		ExposeHeaders:    config.Server.Cors.ExposeHeaders,
		MaxAge:           config.Server.Cors.MaxAge,
	}))

	db := gateways.Connect()

	cert.GenerateKey(logger)
	privateKey, err := oauth2.ReadPrivatekey()
	if err != nil {
		panic(err)
	}

	wconfig := &webauthn.Config{
		RPDisplayName: config.Webauthn.DisplayName,
		RPID:          config.Webauthn.RPID,
		RPOrigins:     config.Webauthn.RPOrigins,
	}

	oauth2 := oauth2.NewOauth2Provider(privateKey, *logger)
	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		panic(err)
	}

	// repositories
	ur := gateways.NewUserRepository(db)
	wcr := gateways.NewWebauthnCredentialRepository(db)
	wsr := gateways.NewWebauthnSessionRepository(db)
	cr := gateways.NewClientRepository(db)
	isr := gateways.NewIdpSessionRepository()
	osr := gateways.NewOidcSessionRepository(db)
	lssr := gateways.NewLoginSkipSessionRepository(db)
	atr := gateways.NewAccessTokenRepository(db)

	// oauth2
	oai := oa.NewAuthorizationInteractor(oauth2, ur, isr, osr, lssr)
	oti := oa.NewTokenInteractor(oauth2, atr)
	oii := oa.NewIntrospectInteractor(oauth2)
	ori := oa.NewRevokeInteractor(oauth2)
	oji := oa.NewJWKInteractor()
	oli := oa.NewLogoutInteractor(oauth2, isr, osr)
	osi := oa.NewSessionInteractor(ur, isr, lssr)
	owi := oa.NewAuthenticateWebauthnInteractor(oauth2, *webAuthn, ur, wcr, lssr)
	controllers.NewOauth2Handler(e, oai, oti, oii, oji, ori, oli, osi, owi)

	// client
	cu := oa.NewClientInteractor(cr)
	controllers.NewClientHandler(e, cu)

	// resource server
	ui := ra.UserinfoInteractor{}
	wi := ra.NewWebauthnInteractor(*logger, *webAuthn, ur, wcr, wsr)
	ii := ra.NewIntrospectInteractor(*logger, atr)
	controllers.NewResourceServerHandler(e, ui, wi, ii)

	e.Logger.Fatal(e.Start(":" + config.Server.Port))
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func loadConfig() (c Config) {
	content, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	var config Config
	if err = yaml.Unmarshal(content, &config); err != nil {
		panic(err)
	}

	return config
}
