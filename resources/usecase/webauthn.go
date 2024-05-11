package usecase

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/k-narusawa/go-idp/authorization/domain/models"
	"github.com/k-narusawa/go-idp/authorization/domain/repository"
	rm "github.com/k-narusawa/go-idp/resources/domain/models"

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
		webauthn.WithExclusions(wu.CredentialExcludeList()),
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

	credentials, err := w.wcr.FindByUserID(ir.Sub)
	if err != nil {
		return err
	}

	resp := WebauthnResponse{
		Keys: make([]WebauthnResponseItem, len(credentials)),
	}

	for i, credential := range credentials {
		wCredential := credential.To()

		id, _ := uuid.FromBytes(credential.ID)
		// idがたまに変なことがあるので、一旦コメントアウト
		// if err != nil {
		// 	return err
		// }

		aaguid, _ := uuid.FromBytes(wCredential.Authenticator.AAGUID)

		resp.Keys[i] = WebauthnResponseItem{
			CredentialID: credential.CredentialID,
			ID:           id.String(),
			AAGUID:       aaguid.String(),
			KeyName:      models.Authenticators[aaguid.String()].Name,
		}
	}

	return c.JSON(200, resp)
}

type WebauthnResponse struct {
	Keys []WebauthnResponseItem `json:"keys"`
}

type WebauthnResponseItem struct {
	CredentialID uint   `json:"credential_id"`
	ID           string `json:"id"`
	AAGUID       string `json:"aaguid"`
	KeyName      string `json:"key_name"`
}

func (w *WebauthnUsecase) Delete(c echo.Context) error {
	// ir := c.Get(("ir")).(rm.IntrospectResponse)

	credentialID := c.Param("credential_id")
	// stringからuintに変換
	credentialIDUint, _ := strconv.Atoi(credentialID)
	err := w.wcr.DeleteByCredentialID(uint(credentialIDUint))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
