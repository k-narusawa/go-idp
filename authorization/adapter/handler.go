package adapter

import (
	"idp/authorization/usecase"

	"github.com/labstack/echo/v4"
)

type Oauth2Handler struct {
	aUsecase usecase.AuthorizationUsecase
	tUseCase usecase.TokenUsecase
	iUseCase usecase.IntrospectUsecase
	jUseCase usecase.JWKUsecase
}

func NewOauth2Handler(e *echo.Echo, authUsecase usecase.AuthorizationUsecase, tokenUsecase usecase.TokenUsecase, introspectUsecase usecase.IntrospectUsecase, jwkUsecase usecase.JWKUsecase) {
	handler := &Oauth2Handler{
		aUsecase: authUsecase,
		tUseCase: tokenUsecase,
		iUseCase: introspectUsecase,
		jUseCase: jwkUsecase,
	}

	e.GET("/oauth2/auth", handler.aUsecase.Invoke)
	e.POST("/oauth2/auth", handler.aUsecase.Invoke)
	e.POST("/oauth2/token", handler.tUseCase.Invoke)
	e.POST("/oauth2/introspect", handler.iUseCase.Invoke)
	e.GET("/oauth2/certs", handler.jUseCase.Invoke)
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
