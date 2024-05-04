package gateway

import (
	"github.com/k-narusawa/go-idp/authorization/domain/models"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type IdSessionRepository struct{}

func NewIdSessionRepository() IdSessionRepository {
	return IdSessionRepository{}
}

func (isr IdSessionRepository) CreateIdSession(c echo.Context, is models.IDSession) error {
	s, err := session.Get("go-idp-session", c)
	if err != nil {
		return err
	}

	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	s.Values["id_session"] = is
	s.Save(c.Request(), c.Response())
	return nil
}

func (isr IdSessionRepository) GetIdSession(c echo.Context) (*models.IDSession, error) {
	s, err := session.Get("go-idp-session", c)
	if err != nil {
		return nil, err
	}

	if s.Values["id_session"] == nil {
		return nil, nil
	}

	idSession := s.Values["id_session"].(*models.IDSession)
	return idSession, nil
}
