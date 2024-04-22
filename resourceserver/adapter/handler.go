package adapter

import (
	"idp/resourceserver/adapter/middleware"
	"idp/resourceserver/usecase"

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

	rs := e.Group("/api/v1/resources")
	rs.Use(middleware.TokenAuthMiddleware())

	rs.GET("/users/userinfo", handler.uu.GetUserinfo)
	rs.GET("/users/webauthn/list", handler.wu.Get)
	rs.GET("/users/webauthn", handler.wu.Start)
	rs.POST("/users/webauthn", handler.wu.Finish)
	rs.DELETE("/users/webauthn", handler.wu.Delete)
}
