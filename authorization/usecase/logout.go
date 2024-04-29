package usecase

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type LogoutUsecase struct {
	oauth2 fosite.OAuth2Provider
}

func NewLogoutUsecase(oauth2 fosite.OAuth2Provider) LogoutUsecase {
	return LogoutUsecase{oauth2: oauth2}
}

func (l *LogoutUsecase) Invoke(c echo.Context) error {
	redirectTo := c.Request().URL.Query().Get("post_logout_redirect_uri")

	sess, _ := session.Get("go-idp-session", c)
	sess.Options = &sessions.Options{MaxAge: -1, Path: "/"}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return c.Redirect(302, redirectTo)
}
