package controllers

import (
	"github.com/k-narusawa/go-idp/authorization/application"

	"github.com/labstack/echo/v4"
)

type Oauth2Handler struct {
	ai application.AuthorizationInteractor
	ti application.TokenInteractor
	ii application.IntrospectInteractor
	ji application.JWKInteractor
	ri application.RevokeInteractor
	li application.LogoutInteractor
	si application.SessionInteractor
	wi application.AuthenticateWebauthnInteractor
}

func NewOauth2Handler(
	e *echo.Echo,
	ai application.AuthorizationInteractor,
	ti application.TokenInteractor,
	ii application.IntrospectInteractor,
	ji application.JWKInteractor,
	ri application.RevokeInteractor,
	li application.LogoutInteractor,
	si application.SessionInteractor,
	awi application.AuthenticateWebauthnInteractor,
) {
	handler := &Oauth2Handler{
		ai: ai,
		ti: ti,
		ii: ii,
		ji: ji,
		ri: ri,
		li: li,
		si: si,
		wi: awi,
	}

	e.GET("/oauth2/auth", handler.ai.Invoke)
	e.POST("/oauth2/auth", handler.ai.Invoke)
	e.POST("/oauth2/token", handler.ti.Invoke)
	e.POST("/oauth2/introspect", handler.ii.Invoke)
	e.POST("/oauth2/revoke", handler.ri.Invoke)
	e.GET("/oauth2/certs", handler.ji.Invoke)
	e.GET("/oauth2/logout", handler.li.Invoke)
	e.GET("/oauth2/session", handler.si.SkipLogin)

	e.GET("/authentication/webauthn/options", handler.wi.Start)
	e.POST("/authentication/webauthn/login", handler.wi.Finish)

	e.GET("/.well-known/openid-configuration", wellKnownOpenIDConfiguration)
}

func wellKnownOpenIDConfiguration(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"issuer":                 "go-idp",
		"authorization_endpoint": "http://localhost:3846/oauth2/auth",
		"token_endpoint":         "http://localhost:3846/oauth2/token",
		"jwks_uri":               "http://localhost:3846/oauth2/certs",
		"response_types_supported": []string{
			"code",
			"token",
			"id_token",
			"code token",
			"code id_token",
			"token id_token",
			"code token id_token",
		},
		"subject_types_supported": []string{
			"public",
		},
		"id_token_signing_alg_values_supported": []string{
			"RS256",
		},
		"scopes_supported": []string{
			"openid",
			"offline",
		},
		"token_endpoint_auth_methods_supported": []string{
			"client_secret_basic",
		},
	})
}
