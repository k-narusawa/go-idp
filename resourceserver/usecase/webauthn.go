package usecase

import (
	"idp/common/adapter/gateway"
	cm "idp/common/domain/models"
	"idp/resourceserver/domain/model"
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
	ir := c.Get(("ir")).(model.IntrospectResponse)
	user := model.NewUser(ir.Sub, "Go-IdP")

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
	}

	options, sd, err := w.webauthn.BeginRegistration(user, registerOptions)
	if err != nil {
		return err
	}

	db := gateway.Connect()

	ws := cm.FromSessionData(sd)

	result := db.Create(&ws)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(http.StatusOK, options)
}

func (w *WebauthnUsecase) Finish(c echo.Context) error {
	ir := c.Get(("ir")).(model.IntrospectResponse)
	user := model.NewUser(ir.Sub, "Go-IdP")

	wsd := cm.WebauthnSessionData{}
	db := gateway.Connect()
	result := db.Where("challenge = ?", c.QueryParam("challenge")).First(&wsd)
	if result.Error != nil {
		return result.Error
	}

	session := wsd.ToSessionData()

	credential, err := w.webauthn.FinishRegistration(user, *session, c.Request())
	if err != nil {
		return err
	}

	user.AddCredential(*credential)

	log.Printf("Credential: %+v\n", credential)

	return c.JSON(http.StatusOK, credential)
}
