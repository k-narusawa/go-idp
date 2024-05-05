package usecase

import (
	"log"
	"net/http"
	"net/url"

	"github.com/k-narusawa/go-idp/authorization/domain/repository"
	"github.com/k-narusawa/go-idp/authorization/oauth2"

	"github.com/k-narusawa/go-idp/authorization/domain/models"

	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type AuthorizationUsecase struct {
	oauth2 fosite.OAuth2Provider
	isr    repository.IIdSessionRepository
	ur     repository.IUserRepository
}

func NewAuthorization(
	oauth2 fosite.OAuth2Provider,
	isr repository.IIdSessionRepository,
	ur repository.IUserRepository,
) AuthorizationUsecase {
	return AuthorizationUsecase{
		oauth2: oauth2,
		isr:    isr,
		ur:     ur,
	}
}

func (a *AuthorizationUsecase) Invoke(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()

	canSkip := false

	is, err := a.isr.GetIdSession(c)
	if err != nil {
		log.Printf("Error occurred in GetIdSession: %+v", err)
		return err
	}

	if is != nil {
		canSkip = true
	}

	if !canSkip {
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
		os := models.NewSession(clientId, user.UserID)

		ar.SetResponseTypeHandled("code")
		response, err := a.oauth2.NewAuthorizeResponse(ctx, ar, os)
		if err != nil {
			log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
			a.oauth2.WriteAuthorizeError(ctx, rw, ar, err)
			return err
		}

		log.Printf("SessionID: %+v", response.GetCode())

		is := models.IDSessionOf(os.Subject, ar)
		is.ClientID = req.PostForm.Get("client_id")

		if err := a.isr.CreateIdSession(c, *is); err != nil {
			log.Printf("Error occurred in CreateIdSession: %+v", err)
			return err
		}

		redirectTo := createRedirectTo(ar, response)

		return c.Redirect(http.StatusFound, redirectTo)
	} else {
		ar := fosite.NewAuthorizeRequest()

		log.Printf("IDSession: %+v", is)

		client, err := oauth2.NewIdpStorage().GetClient(ctx, req.URL.Query().Get("client_id"))
		if err != nil {
			log.Printf("Error occurred in GetClient: %+v", err)
			return err
		}
		ar.Client = client

		redirectURI, _ := url.Parse(req.URL.Query().Get("redirect_uri"))
		ar.RedirectURI = redirectURI

		ar.Form = is.GetRequestForm()
		ar.RequestedAt = is.GetRequestedAt()
		ar.RequestedScope = is.GetRequestedScopes()
		ar.GrantedAudience = is.GetGrantedAudience()
		ar.GrantedScope = is.GetGrantedScopes()
		ar.Session = is.GetSession()
		ar.ID = is.GetID()

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
		response, err := a.oauth2.NewAuthorizeResponse(ctx, ar, is.GetSession())
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
