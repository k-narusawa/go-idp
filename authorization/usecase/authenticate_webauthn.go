package usecase

import (
	"net/http"

	"github.com/k-narusawa/go-idp/authorization/adapter/gateway"
	"github.com/k-narusawa/go-idp/authorization/domain/models"
	"github.com/k-narusawa/go-idp/authorization/domain/repository"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type AuthenticateWebauthnUsecase struct {
	oauth2   fosite.OAuth2Provider
	webauthn webauthn.WebAuthn
	ur       repository.IUserRepository
	wcr      repository.IWebauthnCredentialRepository
	lssr     repository.ILoginSkipSessionRepository
}

func NewAuthenticateWebauthnUsecase(
	oauth2 fosite.OAuth2Provider,
	webauthn webauthn.WebAuthn,
	ur repository.IUserRepository,
	wcr repository.IWebauthnCredentialRepository,
	lssr repository.ILoginSkipSessionRepository,
) AuthenticateWebauthnUsecase {
	return AuthenticateWebauthnUsecase{
		oauth2:   oauth2,
		webauthn: webauthn,
		ur:       ur,
		wcr:      wcr,
		lssr:     lssr,
	}
}

func (w *AuthenticateWebauthnUsecase) Start(c echo.Context) error {
	db := gateway.Connect()

	tx := db.Begin()
	defer tx.Rollback()

	options, sd, err := w.webauthn.BeginDiscoverableLogin(
		webauthn.WithUserVerification(protocol.VerificationRequired),
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

	return c.JSON(200, options.Response)
}

func (w *AuthenticateWebauthnUsecase) Finish(c echo.Context) error {
	sess, err := session.Get("webauthn-session", c)
	if err != nil {
		return err
	}

	sd, ok := sess.Values["authentication"].(*webauthn.SessionData)
	if !ok {
		return c.JSON(http.StatusBadRequest, "session data not found")
	}

	discoverableUserHandler := func(_, userHandle []byte) (webauthn.User, error) {
		user, err := w.ur.FindByUserID(string(userHandle))
		if err != nil {
			return nil, err
		}
		wu := models.NewWebauthnUser(user.UserID, user.Username)
		wc, err := w.wcr.FindByUserID(user.UserID)
		if err != nil {
			return nil, err
		}

		for _, c := range wc {
			wu.AddCredential(*c.To())
		}

		return wu, nil
	}

	parsedResponse, err := protocol.ParseCredentialRequestResponse(c.Request())
	if err != nil {
		return err
	}

	_, err = w.webauthn.ValidateDiscoverableLogin(discoverableUserHandler, *sd, parsedResponse)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{MaxAge: -1, Path: "/"}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	userID := string(parsedResponse.Response.UserHandle)

	lss := models.NewLoginSkipSession(userID)
	err = w.lssr.Save(lss)
	if err != nil {
		return err
	}

	response := WebauthnLoginFinishResponse{
		LoginSkipToken: lss.Token,
	}

	return c.JSON(200, response)
}

type WebauthnLoginFinishResponse struct {
	LoginSkipToken string `json:"login_skip_token"`
}
