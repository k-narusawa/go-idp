package usecase

import (
	"idp/authorization/domain/models"
	"idp/authorization/oauth2"
	"idp/common/adapter/gateway"
	cm "idp/common/domain/models"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
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

	var idSession *models.IDSession
	canSkip := false

	sess, _ := session.Get("go-idp-session", c)
	if sess.Values["id_session"] != nil {
		idSession = sess.Values["id_session"].(*models.IDSession)
		canSkip = true
	}

	if !canSkip {
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

		for _, scope := range req.PostForm["scopes"] {
			ar.GrantScope(scope)
		}

		un := req.PostForm.Get("username")
		p := req.PostForm.Get("password")

		db := gateway.Connect()
		var user cm.User
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

		os := models.NewSession(user.UserID)

		ar.SetResponseTypeHandled("code")
		response, err := a.oauth2.NewAuthorizeResponse(ctx, ar, os)
		if err != nil {
			log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
			a.oauth2.WriteAuthorizeError(ctx, rw, ar, err)
			return err
		}

		is := models.IDSessionOf(os.Subject, ar)
		is.ClientID = req.PostForm.Get("client_id")

		sess, _ := session.Get("go-idp-session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		sess.Values["id_session"] = is
		sess.Save(c.Request(), c.Response())

		redirectTo := createRedirectTo(ar, response)

		return c.Redirect(http.StatusFound, redirectTo)
	} else {
		ar := fosite.NewAuthorizeRequest()

		client, err := oauth2.NewIdpStorage().GetClient(ctx, req.URL.Query().Get("client_id"))
		if err != nil {
			log.Printf("Error occurred in GetClient: %+v", err)
			return err
		}
		ar.Client = client

		redirectURI, _ := url.Parse(req.URL.Query().Get("redirect_uri"))
		ar.RedirectURI = redirectURI

		ar.Form = idSession.GetRequestForm()
		ar.RequestedAt = idSession.GetRequestedAt()
		ar.RequestedScope = idSession.GetRequestedScopes()
		ar.GrantedAudience = idSession.GetGrantedAudience()
		ar.GrantedScope = idSession.GetGrantedScopes()
		ar.Session = idSession.GetSession()
		ar.ID = idSession.GetID()

		ar.ResponseTypes = req.URL.Query()["response_type"]
		ar.State = req.URL.Query().Get("state")

		// nonce対応のため元のリクエストを書き換える
		ar.Form.Del("nonce")
		ar.Form.Add("nonce", req.URL.Query().Get("nonce"))

		// PKCE対応のため元のリクエストを書き換える
		ar.Form.Del("code_challenge")
		ar.Form.Del("code_challenge_method")
		ar.Form.Add("code_challenge", req.URL.Query().Get("code_challenge"))
		ar.Form.Add("code_challenge_method", req.URL.Query().Get("code_challenge_method"))

		ar.SetResponseTypeHandled("code")
		response, err := a.oauth2.NewAuthorizeResponse(ctx, ar, idSession.GetSession())
		if err != nil {
			log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
			a.oauth2.WriteAuthorizeError(ctx, rw, ar, err)
			return err
		}
		redirectTo := createRedirectTo(ar, response)

		return c.Redirect(http.StatusFound, redirectTo)
	}
}

func createRedirectTo(ar fosite.AuthorizeRequester, response fosite.AuthorizeResponder) string {
	redirectTo := ar.GetRedirectURI()

	params := response.GetParameters()

	query := redirectTo.Query()
	for k := range params {
		query.Set(k, params.Get(k))
	}

	redirectTo.RawQuery = query.Encode()

	return redirectTo.String()
}
