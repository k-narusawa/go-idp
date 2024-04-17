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

	if req.Method == "GET" {
		return c.Render(http.StatusOK, "login.html", nil)
	}

	// req.ParseForm()
	// if req.PostForm.Get("username") != "peter" {
	// 	return c.Render(http.StatusOK, "login.html", nil)
	// }

	for _, scope := range req.PostForm["scopes"] {
		ar.GrantScope(scope)
	}

	sess, _ := session.Get("session", c)
	var mySessionData *openid.DefaultSession

	if sess.Values["go-idp"] != nil {
		idpSession := sess.Values["go-idp"].(openid.DefaultSession)
		mySessionData = models.NewSession(idpSession.Subject)
		if err != nil {
			log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
			a.oauth2.WriteAuthorizeError(ctx, rw, ar, err)
			return err
		}
	} else {
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

		mySessionData = models.NewSession(user.UserID)
		sess.Values["go-idp"] = mySessionData
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400,
			HttpOnly: true,
		}

		sess.Save(c.Request(), c.Response())
	}

	// When using the HMACSHA strategy you must use something that implements the HMACSessionContainer.
	// It brings you the power of overriding the default values.
	//
	// mySessionData.HMACSession = &strategy.HMACSession{
	//	AccessTokenExpiry: time.Now().Add(time.Day),
	//	AuthorizeCodeExpiry: time.Now().Add(time.Day),
	// }
	//

	// If you're using the JWT strategy, there's currently no distinction between access token and authorize code claims.
	// Therefore, you both access token and authorize code will have the same "exp" claim. If this is something you
	// need let us know on github.
	//
	// mySessionData.JWTClaims.ExpiresAt = time.Now().Add(time.Day)

	ar.SetResponseTypeHandled("code")
	response, err := a.oauth2.NewAuthorizeResponse(ctx, ar, mySessionData)

	if err != nil {
		log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
		a.oauth2.WriteAuthorizeError(ctx, rw, ar, err)
		return err
	}

	a.oauth2.WriteAuthorizeResponse(ctx, rw, ar, response)

	return nil
}
