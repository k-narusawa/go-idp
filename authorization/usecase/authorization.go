package usecase

import (
	"idp/authorization/adapter/gateway"
	"idp/authorization/domain/models"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
)

type AuthorizationUsecase struct {
	oauth2 fosite.OAuth2Provider
}

func NewAuthorization(oauth2 fosite.OAuth2Provider) AuthorizationUsecase {
	return AuthorizationUsecase{oauth2: oauth2}
}

func (a *AuthorizationUsecase) Invoke(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()

	ar, err := a.oauth2.NewAuthorizeRequest(ctx, req)
	if err != nil {
		log.Printf("Error occurred in NewAuthorizeRequest: %+v", err)
		a.oauth2.WriteAuthorizeError(ctx, rw, ar, err)
		msg := map[string]interface{}{
			"error": "server error",
		}
		return c.Render(http.StatusOK, "login.html", msg)
	}

	var authSession *openid.DefaultSession
	canSkip := false

	idpCookie, _ := c.Cookie("go-idp-session")

	if idpCookie != nil {
		authSession = models.NewSession(idpCookie.Value)

		// TODO: スキップ可能かチェックする
		canSkip = true
	}

	if !canSkip && req.Method == "GET" {
		return c.Render(http.StatusOK, "login.html", nil)
	}

	for _, scope := range req.PostForm["scopes"] {
		ar.GrantScope(scope)
	}

	if !canSkip {
		un := req.PostForm.Get("username")
		p := req.PostForm.Get("password")

		db := gateway.Connect()
		var user models.User
		res := db.Where("username=?", un).First(&user)
		if res.Error != nil {
			log.Printf("Error occurred in GetClient: %+v", res.Error)
			msg := map[string]interface{}{
				"error": "Invalid username or password",
			}
			return c.Render(http.StatusOK, "login.html", msg)
		}

		if err := user.Authenticate(p); err != nil {
			log.Printf("Error occurred in Authenticate: %+v", err)
			msg := map[string]interface{}{
				"error": "Invalid username or password",
			}
			return c.Render(http.StatusOK, "login.html", msg)
		}

		sess, _ := session.Get("go-idp-session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		sess.Values["go-idp-session"] = user.UserID
		sess.Save(c.Request(), c.Response())

		authSession = models.NewSession(user.UserID)
	}

	ar.SetResponseTypeHandled("code")
	response, err := a.oauth2.NewAuthorizeResponse(ctx, ar, authSession)

	if err != nil {
		log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
		a.oauth2.WriteAuthorizeError(ctx, rw, ar, err)
		return err
	}

	a.oauth2.WriteAuthorizeResponse(ctx, rw, ar, response)

	return nil
}
