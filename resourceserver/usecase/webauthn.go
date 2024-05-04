package usecase

import (
	"log"
	"net/http"

	"github.com/k-narusawa/go-idp/authorization/adapter/gateway"
	am "github.com/k-narusawa/go-idp/authorization/domain/models"
	"github.com/k-narusawa/go-idp/resourceserver/domain/models"

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

	u := am.User{}
	result := tx.
		Where("user_id = ?", ir.Sub).
		First(&u)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	wu := am.WebauthnUser{}
	result = tx.Where("name = ?", ir.Sub).First(&wu)

	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			tx.Rollback()
			return result.Error
		}
		wu = *am.NewWebauthnUser(ir.Sub, u.Username)
		result = tx.Create(&wu)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	authSelect := protocol.AuthenticatorSelection{
		// AuthenticatorAttachment: protocol.AuthenticatorAttachment("cross-platform"),
		RequireResidentKey: protocol.ResidentKeyNotRequired(),
		UserVerification:   protocol.VerificationRequired,
	}
	conveyancePref := protocol.PreferNoAttestation

	options, sd, err := w.webauthn.BeginRegistration(wu, webauthn.WithAuthenticatorSelection(authSelect), webauthn.WithConveyancePreference(conveyancePref))

	if err != nil {
		tx.Rollback()
		return err
	}

	ws := am.FromSessionData(sd)

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
	defer tx.Rollback()

	u := am.User{}
	result := tx.
		Where("user_id = ?", ir.Sub).
		First(&u)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	wsd := am.WebauthnSessionData{}
	result = tx.Where("challenge = ?", c.QueryParam("challenge")).First(&wsd)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	session := wsd.ToSessionData()

	wu := am.NewWebauthnUser(u.UserID, u.Username)
	credential, err := w.webauthn.FinishRegistration(wu, *session, c.Request())
	if err != nil {
		return err
	}

	result = tx.Delete(&wsd).Where("challenge = ?", c.QueryParam("challenge"))
	if result.Error != nil {
		log.Printf("Error deleting session data: %+v\n", result.Error)
		return result.Error
	}

	wu.AddCredential(*credential)

	result = tx.Create(&wu.Credentials)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	result = tx.Create(&wu)
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

	wu := am.WebauthnUser{}
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

	wu := am.WebauthnUser{}
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
