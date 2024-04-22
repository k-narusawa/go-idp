package usecase

import (
	"idp/common/adapter/gateway"
	cm "idp/common/domain/models"

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
	db := gateway.Connect()

	tx := db.Begin()
	defer tx.Rollback()

	u := cm.User{}
	result := tx.
		Where("username = ?", "test@example.com").
		First(&u)
	if result.Error != nil {
		return result.Error
	}

	wu := cm.WebauthnUser{}
	result = tx.
		Preload("Credentials").
		Where("id = ?", u.UserID).
		First(&wu)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			return result.Error
		}
	}

	options, session, err := w.webauthn.BeginLogin(wu)
	if err != nil {
		return err
	}

	ws := cm.FromSessionData(session)

	result = tx.Create(&ws)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(200, options)
}
