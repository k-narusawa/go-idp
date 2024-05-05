package adapter

import (
	"github.com/k-narusawa/go-idp/authorization/usecase"

	"github.com/labstack/echo/v4"
)

type Oauth2Handler struct {
	au  usecase.AuthorizationUsecase
	tu  usecase.TokenUsecase
	iu  usecase.IntrospectUsecase
	ju  usecase.JWKUsecase
	ru  usecase.RevokeUsecase
	lu  usecase.LogoutUsecase
	su  usecase.SessionUsecase
	wu  usecase.WebauthnUsecase
	wlu usecase.WebauthnLoginUsecase
}

func NewOauth2Handler(
	e *echo.Echo,
	au usecase.AuthorizationUsecase,
	tu usecase.TokenUsecase,
	iu usecase.IntrospectUsecase,
	ju usecase.JWKUsecase,
	ru usecase.RevokeUsecase,
	lu usecase.LogoutUsecase,
	su usecase.SessionUsecase,
	wl usecase.WebauthnUsecase,
	wlu usecase.WebauthnLoginUsecase,
) {
	handler := &Oauth2Handler{
		au:  au,
		tu:  tu,
		iu:  iu,
		ju:  ju,
		ru:  ru,
		lu:  lu,
		su:  su,
		wu:  wl,
		wlu: wlu,
	}

	e.GET("/oauth2/auth", handler.au.Invoke)
	e.POST("/oauth2/auth", handler.au.Invoke)
	e.POST("/oauth2/token", handler.tu.Invoke)
	e.POST("/oauth2/introspect", handler.iu.Invoke)
	e.POST("/oauth2/revoke", handler.ru.Invoke)
	e.GET("/oauth2/certs", handler.ju.Invoke)
	e.GET("/oauth2/logout", handler.lu.Invoke)
	e.GET("/oauth2/session", handler.lu.Invoke)

	e.GET("/authentication/webauthn/options", handler.wu.Start)
	e.POST("/authentication/webauthn/login", handler.wu.Finish)

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
