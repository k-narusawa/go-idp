package controllers

import (
	"github.com/k-narusawa/go-idp/authorization/application"
	"github.com/k-narusawa/go-idp/middleware"

	"github.com/labstack/echo/v4"
)

type ClientHandler struct {
	ci application.ClientInteractor
}

func NewClientHandler(
	e *echo.Echo,
	ci application.ClientInteractor,
) {
	handler := &ClientHandler{
		ci: ci,
	}

	e.POST("/admin/clients", handler.ci.Register, middleware.InternalAccess())
	e.GET("/admin/clients/:id", handler.ci.Get, middleware.InternalAccess())
}
