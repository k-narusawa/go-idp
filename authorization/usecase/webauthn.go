package usecase

import (
	"idp/common/adapter/gateway"
	"idp/common/domain/models"
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type WebauthnUsecase struct {
	oauth2   fosite.OAuth2Provider
	webauthn webauthn.WebAuthn
}

func NewWebauthnUsecase(oauth2 fosite.OAuth2Provider, webauthn webauthn.WebAuthn) WebauthnUsecase {
	return WebauthnUsecase{oauth2: oauth2, webauthn: webauthn}
}

func (w *WebauthnUsecase) Start(c echo.Context) error {
	db := gateway.Connect()

	tx := db.Begin()
	defer tx.Rollback()

	u := models.User{}
	result := tx.
		Where("username = ?", "test@example.com").
		First(&u)
	if result.Error != nil {
		return result.Error
	}

	wu := models.WebauthnUser{}
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

	options, sd, err := w.webauthn.BeginLogin(
		wu,
		webauthn.WithAllowedCredentials(allowList),
		webauthn.WithUserVerification(protocol.VerificationRequired),
		webauthn.WithAppIdExtension(""),
	)
	if err != nil {
		return err
	}

	sess, err := session.Get("webauthn-session", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	sess.Values["authentication"] = sd
	sess.Save(c.Request(), c.Response())

	tx.Commit()

	return c.JSON(200, options)
}

func (w *WebauthnUsecase) Finish(c echo.Context) error {
	sess, err := session.Get("webauthn-session", c)
	if err != nil {
		return err
	}

	sd, ok := sess.Values["authentication"].(*webauthn.SessionData)
	if !ok {
		return c.JSON(http.StatusBadRequest, "session data not found")
	}

	db := gateway.Connect()

	tx := db.Begin()
	defer tx.Rollback()

	u := models.User{}
	result := tx.
		Where("username = ?", "test@example.com").
		First(&u)
	if result.Error != nil {
		return result.Error
	}

	wu := models.WebauthnUser{}
	result = tx.
		Preload("Credentials").
		Where("id = ?", u.UserID).
		First(&wu)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Error finding user: %+v\n", result.Error)
		return result.Error
	}

	_, err = w.webauthn.FinishLogin(wu, *sd, c.Request())
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{MaxAge: -1, Path: "/"}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	tx.Commit()

	return nil
}