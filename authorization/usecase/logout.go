package usecase

import (
	"github.com/gorilla/sessions"
	"github.com/k-narusawa/go-idp/domain/repository"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type LogoutUsecase struct {
	oauth2 fosite.OAuth2Provider
	isr    repository.IIdpSessionRepository
	osr    repository.IOidcSessionRepository
}

func NewLogoutUsecase(
	oauth2 fosite.OAuth2Provider,
	isr repository.IIdpSessionRepository,
	osr repository.IOidcSessionRepository,
) LogoutUsecase {
	return LogoutUsecase{
		oauth2: oauth2,
		isr:    isr,
		osr:    osr,
	}
}

func (l *LogoutUsecase) Invoke(c echo.Context) error {
	redirectTo := c.Request().URL.Query().Get("post_logout_redirect_uri")

	// DBからSessionを削除
	oidcSession, err := l.isr.Get(c)
	if err != nil {
		return err
	}

	if oidcSession != nil {
		l.osr.DeleteBySignature(oidcSession.SessionID)
	}

	// Cookie上のSessionを全て削除
	sess, _ := session.Get("go-idp-session", c)
	sess.Options = &sessions.Options{MaxAge: -1, Path: "/"}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return c.Redirect(302, redirectTo)
}
