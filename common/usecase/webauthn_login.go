package usecase

import (
	"idp/common/adapter/gateway"
	cm "idp/common/domain/models"
	"log"

	"github.com/go-webauthn/webauthn/protocol"
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

	allowList := make([]protocol.CredentialDescriptor, len(wu.Credentials))
	for i := range wu.Credentials {
		wc := wu.Credentials[i].ToWebauthnCredential()

		allowList[i] = protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: wc.Descriptor().CredentialID,
			Transport:    []protocol.AuthenticatorTransport{"usb", "internal", "hybrid", "ble", "nfc"},
		}
	}

	options, session, err := w.webauthn.BeginLogin(
		wu,
		webauthn.WithAllowedCredentials(allowList),
		webauthn.WithUserVerification(protocol.VerificationRequired),
		webauthn.WithAppIdExtension(""),
	)
	if err != nil {
		return err
	}

	ws := cm.FromSessionData(session)

	result = tx.Create(&ws)
	if result.Error != nil {
		return result.Error
	}

	tx.Commit()

	return c.JSON(200, options)
}

func (w *WebauthnLoginUsecase) Finish(c echo.Context) error {
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

	wsd := cm.WebauthnSessionData{}
	result = tx.Where("challenge = ?", c.QueryParam("challenge")).First(&wsd)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Error finding session data: %+v\n", result.Error)
		return result.Error
	}

	session := wsd.ToSessionData()

	wu := cm.WebauthnUser{}
	result = tx.
		Preload("Credentials").
		Where("id = ?", u.UserID).
		First(&wu)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Error finding user: %+v\n", result.Error)
		return result.Error
	}

	_, err := w.webauthn.FinishLogin(wu, *session, c.Request())
	if err != nil {
		return err
	}

	result = tx.Delete(&wsd).Where("challenge = ?", c.QueryParam("challenge"))
	if result.Error != nil {
		log.Printf("Error deleting session data: %+v\n", result.Error)
		return result.Error
	}

	tx.Commit()

	return nil
}
