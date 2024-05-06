package usecase

import (
	"net/http"

	"github.com/google/uuid"
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
	wsr      repository.IWebauthnSessionRepository
}

func NewWebauthnUsecase(
	webauthn webauthn.WebAuthn,
	ur repository.IUserRepository,
	wcr repository.IWebauthnCredentialRepository,
	wsr repository.IWebauthnSessionRepository,
) WebauthnUsecase {
	return WebauthnUsecase{
		webauthn: webauthn,
		ur:       ur,
		wcr:      wcr,
		wsr:      wsr,
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

	err = w.wsr.Save(ws)
	if err != nil {
		return err
	}

	return c.JSON(200, options.Response)
}

func (w *WebauthnUsecase) Finish(c echo.Context) error {
	ir := c.Get(("ir")).(rm.IntrospectResponse)
	challenge := c.QueryParam("challenge")

	user, err := w.ur.FindByUserID(ir.Sub)
	if err != nil {
		return err
	}

	wsd, err := w.wsr.FindByChallenge(challenge)
	if err != nil {
		return err
	}

	session := wsd.ToSessionData()

	wu := models.NewWebauthnUser(user.UserID, user.Username)

	credential, err := w.webauthn.FinishRegistration(wu, *session, c.Request())
	if err != nil {
		return err
	}

	err = w.wsr.DeleteByChallenge(challenge)
	if err != nil {
		return err
	}

	if w.wcr.Save(models.FromWebauthnCredential(user.UserID, credential)) != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (w *WebauthnUsecase) Get(c echo.Context) error {
	ir := c.Get(("ir")).(rm.IntrospectResponse)

	wcs, err := w.wcr.FindByUserID(ir.Sub)
	if err != nil {
		return err
	}

	credentials := make([]webauthn.Credential, len(wcs))
	for i, wc := range wcs {
		credentials[i] = *wc.To()
	}

	resp := WebauthnResponse{
		Keys: make([]WebauthnResponseItem, len(credentials)),
	}

	for i, cred := range credentials {
		id, err := uuid.FromBytes(cred.ID)
		if err != nil {
			return err
		}

		aaguid, _ := uuid.FromBytes(cred.Authenticator.AAGUID)

		resp.Keys[i] = WebauthnResponseItem{
			ID:      id.String(),
			AAGUID:  aaguid.String(),
			KeyName: models.Authenticators[aaguid.String()].Name,
		}
	}

	return c.JSON(200, resp)
}

type WebauthnResponse struct {
	Keys []WebauthnResponseItem `json:"keys"`
}

type WebauthnResponseItem struct {
	ID      string `json:"id"`
	AAGUID  string `json:"aaguid"`
	KeyName string `json:"key_name"`
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
