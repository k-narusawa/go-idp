package adapter

import (
	"idp/resourceserver/usecase"

	"github.com/labstack/echo/v4"
)

type ResourceServerHandler struct {
	uu usecase.UserinfoUsecase
}

func NewResourceServerHandler(e *echo.Echo, uu usecase.UserinfoUsecase) {
	handler := &ResourceServerHandler{
		uu: uu,
	}
	e.GET("/api/userinfo", handler.uu.GetUserinfo)
}
