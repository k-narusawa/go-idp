package adapter

import (
	"github.com/k-narusawa/go-idp/resourceserver/adapter/middleware"
	"github.com/k-narusawa/go-idp/resourceserver/usecase"

	"github.com/labstack/echo/v4"
)

type ResourceServerHandler struct {
	uu usecase.UserinfoUsecase
	wu usecase.WebauthnUsecase
}

func NewResourceServerHandler(e *echo.Echo, uu usecase.UserinfoUsecase, wu usecase.WebauthnUsecase) {
	handler := &ResourceServerHandler{
		uu: uu,
		wu: wu,
	}

	r := e.Group("/resources")
	r.Use(middleware.TokenAuthMiddleware())

	r.GET("/users/userinfo", handler.uu.GetUserinfo)
	// rs.GET("/users/webauthn/list", handler.wu.Get)
	r.GET("/users/registrations/webauthn/options", handler.wu.Start)
	r.POST("/users/registrations/webauthn/result", handler.wu.Finish)
	// r.DELETE("/users/webauthn", handler.wu.Delete)
}
