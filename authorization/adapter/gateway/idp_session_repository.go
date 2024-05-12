package gateway

import (
	"github.com/gorilla/sessions"
	"github.com/k-narusawa/go-idp/domain/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type IdpSessionRepository struct{}

func NewIdpSessionRepository() IdpSessionRepository {
	return IdpSessionRepository{}
}

func (isr IdpSessionRepository) Save(c echo.Context, idpSession *models.IdpSession) error {
	s, err := session.Get("go-idp-session", c)
	if err != nil {
		return err
	}

	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	s.Values["idp_session"] = idpSession
	s.Save(c.Request(), c.Response())
	return nil
}

func (isr IdpSessionRepository) Get(c echo.Context) (*models.IdpSession, error) {
	s, err := session.Get("go-idp-session", c)
	if err != nil {
		return nil, err
	}

	if s.Values["idp_session"] == nil {
		return nil, nil
	}

	idpSession := s.Values["idp_session"].(*models.IdpSession)
	return idpSession, nil
}

func (isr IdpSessionRepository) Delete(c echo.Context) error {
	s, err := session.Get("go-idp-session", c)
	if err != nil {
		return err
	}

	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	s.Values["idp_session"] = nil
	s.Save(c.Request(), c.Response())
	return nil
}
