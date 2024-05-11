package adapter

import (
	"github.com/k-narusawa/go-idp/middleware"
	"github.com/k-narusawa/go-idp/resources/usecase"

	"github.com/labstack/echo/v4"
)

type ResourceServerHandler struct {
	uu usecase.UserinfoUsecase
	wu usecase.WebauthnUsecase
	iu usecase.IntrospectUsecase
}

func NewResourceServerHandler(
	e *echo.Echo,
	uu usecase.UserinfoUsecase,
	wu usecase.WebauthnUsecase,
	iu usecase.IntrospectUsecase,
) {
	handler := &ResourceServerHandler{
		uu: uu,
		wu: wu,
		iu: iu,
	}

	r := e.Group("/resources")
	r.Use(middleware.TokenAuthMiddleware(iu))

	r.GET("/users/userinfo", handler.uu.GetUserinfo)
	r.GET("/users/registrations/webauthn/options", handler.wu.Start)
	r.POST("/users/registrations/webauthn/result", handler.wu.Finish)
	r.GET("/users/webauthn", handler.wu.Get)
	r.DELETE("/users/webauthn/:credential_id", handler.wu.Delete)
}
