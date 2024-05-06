package usecase

import (
	"log"
	"net/http"

	"github.com/k-narusawa/go-idp/authorization/adapter/gateway"
	"github.com/k-narusawa/go-idp/authorization/domain/models"
	"github.com/k-narusawa/go-idp/authorization/domain/repository"
	rm "github.com/k-narusawa/go-idp/resourceserver/domain/models"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

type WebauthnUsecase struct {
	webauthn webauthn.WebAuthn
	ur       repository.IUserRepository
	wcr      repository.IWebauthnCredentialRepository
}

func NewWebauthnUsecase(
	webauthn webauthn.WebAuthn,
	ur repository.IUserRepository,
	wcr repository.IWebauthnCredentialRepository,
) WebauthnUsecase {
	return WebauthnUsecase{
		webauthn: webauthn,
		ur:       ur,
		wcr:      wcr,
	}
}

func (w *WebauthnUsecase) Start(c echo.Context) error {
	ir := c.Get(("ir")).(rm.IntrospectResponse)

	user, err := w.ur.FindByUserID(ir.Sub)
	if err != nil {
		return err
	}

	wu := models.NewWebauthnUser(user.UserID, user.Username)
	credentials, err := w.wcr.FindByUserID(user.UserID)
	if err != nil {
		return err
	}

	for _, cred := range credentials {
		wu.AddCredential(*cred.To())
	}

	authSelect := protocol.AuthenticatorSelection{
		// AuthenticatorAttachment: protocol.AuthenticatorAttachment("cross-platform"),
		RequireResidentKey: protocol.ResidentKeyNotRequired(),
		UserVerification:   protocol.VerificationRequired,
	}
	conveyancePref := protocol.PreferNoAttestation

	options, session, err := w.webauthn.BeginRegistration(
		wu,
		webauthn.WithAuthenticatorSelection(authSelect),
		webauthn.WithConveyancePreference(conveyancePref),
	)

	if err != nil {
		return err
	}

	ws := models.FromSessionData(session)

	db := gateway.Connect()
	result := db.Create(&ws)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(200, options.Response)
}

func (w *WebauthnUsecase) Finish(c echo.Context) error {
	ir := c.Get(("ir")).(rm.IntrospectResponse)
	db := gateway.Connect()

	user, err := w.ur.FindByUserID(ir.Sub)
	if err != nil {
		return err
	}

	wsd := models.WebauthnSessionData{}
	result := db.Where("challenge = ?", c.QueryParam("challenge")).First(&wsd)
	if result.Error != nil {
		db.Rollback()
		log.Printf("Error finding session data: %+v\n", result.Error)
		return result.Error
	}

	session := wsd.ToSessionData()

	wu := models.NewWebauthnUser(user.UserID, user.Username)

	credential, err := w.webauthn.FinishRegistration(wu, *session, c.Request())
	if err != nil {
		return err
	}

	result = db.Delete(&wsd).Where("challenge = ?", c.QueryParam("challenge"))
	if result.Error != nil {
		log.Printf("Error deleting session data: %+v\n", result.Error)
		return result.Error
	}

	w.wcr.Save(models.FromWebauthnCredential(user.UserID, credential))

	return nil
}

func (w *WebauthnUsecase) Delete(c echo.Context) error {
	ir := c.Get(("ir")).(rm.IntrospectResponse)

	db := gateway.Connect()
	tx := db.Begin()
	defer tx.Rollback()

	wu := models.WebauthnUser{}
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
