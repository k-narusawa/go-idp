package adapter

import (
	"idp/common/usecase"

	"github.com/labstack/echo/v4"
)

type CommonHandler struct {
	wl usecase.WebauthnLoginUsecase
}

func NewCommonHandler(e *echo.Echo, wl usecase.WebauthnLoginUsecase) {
	handler := &CommonHandler{
		wl: wl,
	}

	cm := e.Group("api/v1")
	cm.GET("/webauthn/login", handler.wl.Start)
	cm.POST("/webauthn/login", handler.wl.Finish)
}
