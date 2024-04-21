package usecase

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

type WebauthnLoginUsecase struct {
	webauthn webauthn.WebAuthn
}

func NewWebauthnLoginUsecase(webauthn webauthn.WebAuthn) WebauthnLoginUsecase {
	return WebauthnLoginUsecase{webauthn: webauthn}
}

func (w *WebauthnLoginUsecase) Start(c echo.Context) error {
	return c.JSON(200, "Hello, World!")
}
