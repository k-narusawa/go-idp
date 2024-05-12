package usecase

import (
	"log"

	"github.com/k-narusawa/go-idp/domain/models"
	"github.com/k-narusawa/go-idp/domain/repository"

	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type TokenUsecase struct {
	oauth2 fosite.OAuth2Provider
	atr    repository.IAccessTokenRepository
}

func NewTokenUsecase(
	oauth2 fosite.OAuth2Provider,
	atr repository.IAccessTokenRepository,
) TokenUsecase {
	return TokenUsecase{
		oauth2: oauth2,
		atr:    atr,
	}
}

func (t *TokenUsecase) Invoke(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()

	// Create an empty session object which will be passed to the request handlers
	emptySession := models.NewEmptyIdpSession()

	// This will create an access request object and iterate through the registered TokenEndpointHandlers to validate the request.
	ar, err := t.oauth2.NewAccessRequest(ctx, req, emptySession)

	if err != nil {
		log.Printf("Error occurred in NewAccessRequest: %+v", err)
		t.oauth2.WriteAccessError(ctx, rw, ar, err)
		return nil
	}

	// If this is a client_credentials grant, grant all requested scopes
	// NewAccessRequest validated that all requested scopes the client is allowed to perform
	// based on configured scope matching strategy.
	if ar.GetGrantTypes().ExactOne("client_credentials") {
		for _, scope := range ar.GetRequestedScopes() {
			ar.GrantScope(scope)
		}
	}

	// Next we create a response for the access request. Again, we iterate through the TokenEndpointHandlers
	// and aggregate the result in response.
	response, err := t.oauth2.NewAccessResponse(ctx, ar)
	if err != nil {
		log.Printf("Error occurred in NewAccessResponse: %+v", err)
		t.oauth2.WriteAccessError(ctx, rw, ar, err)
		return nil
	}

	subject := ar.GetSession().GetSubject()
	accessTokenList, err := t.atr.FindBySubject(subject)
	if err != nil {
		log.Printf("Error occurred in FindBySubject: %+v", err)
	}

	for _, accessToken := range *accessTokenList {
		if !accessToken.IsExpired() {
			continue
		}
		err := t.atr.DeleteBySignature(accessToken.Signature)
		if err != nil {
			log.Printf("Error occurred in DeleteBySignature: %+v", err)
		}
	}

	t.oauth2.WriteAccessResponse(ctx, rw, ar, response)

	return nil

}
