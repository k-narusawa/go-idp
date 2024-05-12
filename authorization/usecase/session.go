package usecase

import (
	"net/http"
	"net/url"

	"github.com/k-narusawa/go-idp/domain/models"
	"github.com/k-narusawa/go-idp/domain/repository"
	"github.com/labstack/echo/v4"
)

type SessionUsecase struct {
	ur   repository.IUserRepository
	isr  repository.IIdpSessionRepository
	lssr repository.ILoginSkipSessionRepository
}

func NewSessionUsecase(
	ur repository.IUserRepository,
	isr repository.IIdpSessionRepository,
	lssr repository.ILoginSkipSessionRepository,
) SessionUsecase {
	return SessionUsecase{
		ur:   ur,
		isr:  isr,
		lssr: lssr,
	}
}

func (s *SessionUsecase) SkipLogin(c echo.Context) error {
	token := c.QueryParam("token")
	clientId := c.QueryParam("client_id")
	scope := c.QueryParam("scope")
	responseType := c.QueryParam("response_type")
	redirectUri := c.QueryParam("redirect_uri")
	grantType := c.QueryParam("grant_type")
	state := c.QueryParam("state")
	codeChallenge := c.QueryParam("code_challenge")
	codeChallengeMethod := c.QueryParam("code_challenge_method")
	nonce := c.QueryParam("nonce")

	lss, err := s.lssr.FindByToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid token"})
	}

	if lss.IsExpired() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "token expired"})
	}

	redirectTo, _ := url.Parse("/oauth2/auth")
	params := url.Values{}
	params.Add("client_id", clientId)
	params.Add("scope", scope)
	params.Add("response_type", responseType)
	params.Add("redirect_uri", redirectUri)
	params.Add("grant_type", grantType)
	params.Add("state", state)
	if codeChallenge != "" && codeChallengeMethod != "" {
		params.Add("code_challenge", codeChallenge)
		params.Add("code_challenge_method", codeChallengeMethod)
	}

	if nonce != "" {
		params.Add("nonce", nonce)
	}

	user, err := s.ur.FindByUserID(lss.UserID)
	if err != nil {
		c.Redirect(http.StatusFound, "/error")
	}

	redirectTo.RawQuery = params.Encode()
	idpSession := models.NewIdpSession(clientId, *user)
	idpSession.SetLoginSkipToken(token)

	s.isr.Save(c, idpSession)

	return c.Redirect(http.StatusFound, redirectTo.String())
}
