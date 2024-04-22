package usecase

import (
	"idp/common/adapter/gateway"
	cm "idp/common/domain/models"
	"idp/resourceserver/domain/models"
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

type WebauthnUsecase struct {
	webauthn webauthn.WebAuthn
}

func NewWebauthnUsecase(webauthn webauthn.WebAuthn) WebauthnUsecase {
	return WebauthnUsecase{webauthn: webauthn}
}

func (w *WebauthnUsecase) Start(c echo.Context) error {
	ir := c.Get(("ir")).(models.IntrospectResponse)

	db := gateway.Connect()
	tx := db.Begin()
	defer tx.Rollback()

	user := cm.WebauthnUser{}
	result := tx.Where("name = ?", ir.Sub).First(&user)

	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			tx.Rollback()
			return result.Error
		}
		user = *cm.NewWebauthnUser(ir.Sub, "Go-IdP")
		result = tx.Create(&user)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
	}

	options, sd, err := w.webauthn.BeginRegistration(user, registerOptions)
	if err != nil {
		tx.Rollback()
		return err
	}

	ws := cm.FromSessionData(sd)

	result = tx.Create(&ws)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()

	return c.JSON(http.StatusOK, options)
}

func (w *WebauthnUsecase) Finish(c echo.Context) error {
	ir := c.Get(("ir")).(models.IntrospectResponse)

	db := gateway.Connect()
	tx := db.Begin()

	user := cm.NewWebauthnUser(ir.Sub, "Go-IdP")

	wsd := cm.WebauthnSessionData{}
	result := tx.Where("challenge = ?", c.QueryParam("challenge")).First(&wsd)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	session := wsd.ToSessionData()

	credential, err := w.webauthn.FinishRegistration(user, *session, c.Request())
	if err != nil {
		return err
	}

	result = tx.Delete(&wsd).Where("challenge = ?", c.QueryParam("challenge"))
	if result.Error != nil {
		log.Printf("Error deleting session data: %+v\n", result.Error)
		return result.Error
	}

	user.AddCredential(*credential)

	result = tx.Create(&user.Credentials)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	result = tx.Create(&user)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()

	return c.JSON(http.StatusOK, credential)
}

func (w *WebauthnUsecase) Get(c echo.Context) error {
	ir := c.Get(("ir")).(models.IntrospectResponse)

	db := gateway.Connect()
	tx := db.Begin()

	wu := cm.WebauthnUser{}
	result := tx.
		Preload("Credentials").
		Where("name = ?", ir.Sub).
		Find(&wu)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	cred := wu.WebAuthnCredentials()

	tx.Commit()

	return c.JSON(http.StatusOK, cred)
}

func (w *WebauthnUsecase) Delete(c echo.Context) error {
	ir := c.Get(("ir")).(models.IntrospectResponse)

	db := gateway.Connect()
	tx := db.Begin()
	defer tx.Rollback()

	wu := cm.WebauthnUser{}
	result := tx.
		Preload("Credentials").
		Where("name = ?", ir.Sub).
		First(&wu)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Delete(&wu.Credentials)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()

	return c.NoContent(http.StatusNoContent)
}
