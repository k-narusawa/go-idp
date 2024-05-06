package usecase

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/k-narusawa/go-idp/authorization/domain/repository"

	"github.com/k-narusawa/go-idp/authorization/domain/models"

	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type AuthorizationUsecase struct {
	oauth2 fosite.OAuth2Provider
	ur     repository.IUserRepository
	isr    repository.IIdpSessionRepository
	osr    repository.IOidcSessionRepository
	lssr   repository.ILoginSkipSessionRepository
}

func NewAuthorization(
	oauth2 fosite.OAuth2Provider,
	ur repository.IUserRepository,
	isr repository.IIdpSessionRepository,
	osr repository.IOidcSessionRepository,
	lssr repository.ILoginSkipSessionRepository,
) AuthorizationUsecase {
	return AuthorizationUsecase{
		oauth2: oauth2,
		ur:     ur,
		isr:    isr,
		osr:    osr,
		lssr:   lssr,
	}
}

func (a *AuthorizationUsecase) Invoke(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()

	canSkip := false
	hasLoginSkipToken := false

	idpSession, _ := a.isr.Get(c)

	if idpSession != nil {
		if idpSession.LoginSkipToken == "" {
			canSkip = true
		} else if idpSession.LoginSkipToken != "" {
			hasLoginSkipToken = true
		}
	}

	if !canSkip && !hasLoginSkipToken {
		ar, err := a.oauth2.NewAuthorizeRequest(ctx, req)
		if err != nil {
			log.Printf("Error occurred in NewAuthorizeRequest: %+v", err)
			a.oauth2.WriteAuthorizeError(ctx, rw, ar, err)
			msg := "username or password is invalid."
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

		user, err := a.ur.FindByUsername(un)
		if err != nil {
			msg := "username or password is invalid."
			return c.Render(http.StatusOK, "login.html", msg)
		}

		if err := user.Authenticate(p); err != nil {
			log.Printf("Error occurred in Authenticate: %+v", err)
			msg := "username or password is invalid."
			return c.Render(http.StatusOK, "login.html", msg)
		}

		clientId := ar.GetClient().GetID()
		idpSession := models.NewIdpSession(clientId, *user)

		ar.SetResponseTypeHandled("code")
		response, err := a.oauth2.NewAuthorizeResponse(ctx, ar, idpSession)
		if err != nil {
			log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
			a.oauth2.WriteAuthorizeError(ctx, rw, ar, err)
			return err
		}

		idpSession.SetSessionID(response.GetCode())

		if err := a.isr.Save(c, idpSession); err != nil {
			log.Printf("Error occurred in CreateIdSession: %+v", err)
			return err
		}

		redirectTo := createRedirectTo(ar, response)

		return c.Redirect(http.StatusFound, redirectTo)
	} else if hasLoginSkipToken {
		ar, err := a.oauth2.NewAuthorizeRequest(ctx, req)
		if err != nil {
			log.Printf("Error occurred in NewAuthorizeRequest: %+v", err)
			a.oauth2.WriteAuthorizeError(ctx, rw, ar, err)
			msg := "username or password is invalid."
			return c.Render(http.StatusOK, "login.html", msg)
		}

		scopes := strings.Split(req.URL.Query()["scope"][0], " ")
		log.Printf("scope: %+v", scopes)
		for _, scope := range scopes {
			ar.GrantScope(scope)
		}

		lss, err := a.lssr.FindByToken(idpSession.LoginSkipToken)
		if err != nil {
			msg := "unexpected error occurred."
			return c.Render(http.StatusOK, "login.html", msg)
		}

		clientId := ar.GetClient().GetID()
		user, err := a.ur.FindByUserID(lss.UserID)
		if err != nil {
			msg := "unexpected error occurred."
			return c.Render(http.StatusOK, "login.html", msg)
		}
		idpSession := models.NewIdpSession(clientId, *user)

		ar.SetResponseTypeHandled("code")
		response, err := a.oauth2.NewAuthorizeResponse(ctx, ar, idpSession)
		if err != nil {
			log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
			a.oauth2.WriteAuthorizeError(ctx, rw, ar, err)
			return err
		}

		idpSession.SetSessionID(response.GetCode())
		idpSession.RemoveLoginSkipToken()

		if err := a.isr.Save(c, idpSession); err != nil {
			msg := "unexpected error occurred."
			return c.Render(http.StatusOK, "login.html", msg)
		}

		redirectTo := createRedirectTo(ar, response)

		return c.Redirect(http.StatusFound, redirectTo)
	} else {
		ar := fosite.NewAuthorizeRequest()

		oidcSession, err := a.osr.FindBySignature(idpSession.SessionID)
		if err != nil {
			log.Printf("Error occurred in FindBySignature: %+v", err)
			return err
		}

		redirectURI, _ := url.Parse(req.URL.Query().Get("redirect_uri"))
		ar.RedirectURI = redirectURI

		ar.Form = oidcSession.GetRequestForm()
		ar.RequestedAt = oidcSession.GetRequestedAt()
		ar.RequestedScope = oidcSession.GetRequestedScopes()
		ar.GrantedAudience = oidcSession.GetGrantedAudience()
		ar.GrantedScope = oidcSession.GetGrantedScopes()
		ar.Session = oidcSession.GetSession()
		ar.ID = oidcSession.GetID()
		ar.Client = oidcSession.GetClient()

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
		response, err := a.oauth2.NewAuthorizeResponse(ctx, ar, oidcSession.GetSession())
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
