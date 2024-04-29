package usecase

import (
	"log"

	"github.com/k-narusawa/go-idp/authorization/domain/models"

	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type TokenUsecase struct {
	oauth2 fosite.OAuth2Provider
}

func NewTokenUsecase(oauth2 fosite.OAuth2Provider) TokenUsecase {
	return TokenUsecase{oauth2: oauth2}
}

func (t *TokenUsecase) Invoke(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()

	// Create an empty session object which will be passed to the request handlers
	mySessionData := models.NewSession("")

	// This will create an access request object and iterate through the registered TokenEndpointHandlers to validate the request.
	accessRequest, err := t.oauth2.NewAccessRequest(ctx, req, mySessionData)

	if err != nil {
		log.Printf("Error occurred in NewAccessRequest: %+v", err)
		t.oauth2.WriteAccessError(ctx, rw, accessRequest, err)
		return nil
	}

	// If this is a client_credentials grant, grant all requested scopes
	// NewAccessRequest validated that all requested scopes the client is allowed to perform
	// based on configured scope matching strategy.
	if accessRequest.GetGrantTypes().ExactOne("client_credentials") {
		for _, scope := range accessRequest.GetRequestedScopes() {
			accessRequest.GrantScope(scope)
		}
	}

	// Next we create a response for the access request. Again, we iterate through the TokenEndpointHandlers
	// and aggregate the result in response.
	response, err := t.oauth2.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		log.Printf("Error occurred in NewAccessResponse: %+v", err)
		t.oauth2.WriteAccessError(ctx, rw, accessRequest, err)
		return nil
	}

	// All done, send the response.
	t.oauth2.WriteAccessResponse(ctx, rw, accessRequest, response)

	return nil

}
