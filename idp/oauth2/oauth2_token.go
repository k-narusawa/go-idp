package oauth2

import (
	"log"

	"github.com/labstack/echo/v4"
)

func TokenEndpoint(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()

	// Create an empty session object which will be passed to the request handlers
	mySessionData := newSession("")

	// This will create an access request object and iterate through the registered TokenEndpointHandlers to validate the request.
	accessRequest, err := oauth2.NewAccessRequest(ctx, req, mySessionData)

	if err != nil {
		log.Printf("Error occurred in NewAccessRequest: %+v", err)
		oauth2.WriteAccessError(rw, accessRequest, err)
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
	response, err := oauth2.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		log.Printf("Error occurred in NewAccessResponse: %+v", err)
		oauth2.WriteAccessError(rw, accessRequest, err)
		return nil
	}

	// All done, send the response.
	oauth2.WriteAccessResponse(rw, accessRequest, response)

	return nil
}