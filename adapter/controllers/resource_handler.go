package controllers

import (
	"github.com/k-narusawa/go-idp/middleware"
	"github.com/k-narusawa/go-idp/resources/application"

	"github.com/labstack/echo/v4"
)

type ResourceServerHandler struct {
	uu application.UserinfoInteractor
	wu application.WebauthnInteractor
	iu application.IntrospectInteractor
}

func NewResourceServerHandler(
	e *echo.Echo,
	ui application.UserinfoInteractor,
	wi application.WebauthnInteractor,
	ii application.IntrospectInteractor,
) {
	handler := &ResourceServerHandler{
		uu: ui,
		wu: wi,
		iu: ii,
	}

	r := e.Group("/resources")
	r.Use(middleware.TokenAuthMiddleware(ii))

	r.GET("/users/userinfo", handler.uu.GetUserinfo)
	r.GET("/users/registrations/webauthn/options", handler.wu.Start)
	r.POST("/users/registrations/webauthn/result", handler.wu.Finish)
	r.GET("/users/webauthn", handler.wu.Get)
	r.DELETE("/users/webauthn/:credential_id", handler.wu.Delete)
}
