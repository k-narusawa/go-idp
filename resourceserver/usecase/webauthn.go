package usecase

import (
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

	options, session, err := w.webauthn.BeginRegistration(user, registerOptions)
	if err != nil {
		return err
	}

	log.Printf("Session: %+v\n", session)

	return c.JSON(http.StatusOK, options)
}
