package adapter

import (
	"idp/authorization/usecase"

	"github.com/labstack/echo/v4"
)

type Oauth2Handler struct {
	aUsecase usecase.AuthorizationUsecase
	tUseCase usecase.TokenUsecase
	iUseCase usecase.IntrospectUsecase
}

func NewOauth2Handler(e *echo.Echo, authUsecase usecase.AuthorizationUsecase, tokenUsecase usecase.TokenUsecase, introspectUsecase usecase.IntrospectUsecase) {
	handler := &Oauth2Handler{
		aUsecase: authUsecase,
		tUseCase: tokenUsecase,
		iUseCase: introspectUsecase,
	}

	e.GET("/oauth2/auth", handler.aUsecase.Invoke)
	e.POST("/oauth2/auth", handler.aUsecase.Invoke)
	e.POST("/oauth2/token", handler.tUseCase.Invoke)
	e.POST("/oauth2/introspect", handler.iUseCase.Invoke)
}
