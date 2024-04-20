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

	rs := e.Group("/api/v1/users")
	rs.Use(middleware.TokenAuthMiddleware())

	rs.GET("/userinfo", handler.uu.GetUserinfo)
	rs.GET("/webauthn", handler.wu.Start)
	rs.POST("/webauthn", handler.wu.Finish)
}
