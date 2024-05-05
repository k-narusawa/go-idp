package adapter

import (
	"github.com/k-narusawa/go-idp/authorization/adapter/middleware"
	"github.com/k-narusawa/go-idp/authorization/usecase"

	"github.com/labstack/echo/v4"
)

type ClientHandler struct {
	cu usecase.ClientUsecase
}

func NewClientHandler(
	e *echo.Echo,
	cu usecase.ClientUsecase,
) {
	handler := &ClientHandler{
		cu: cu,
	}

	e.POST("/admin/clients", handler.cu.Register, middleware.InternalAccess())
	e.GET("/admin/clients/:id", handler.cu.Get, middleware.InternalAccess())
}
