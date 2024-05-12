package application

import (
	"log"

	"github.com/k-narusawa/go-idp/domain/models"

	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type IntrospectInteractor struct {
	oauth2 fosite.OAuth2Provider
}

func NewIntrospectInteractor(oauth2 fosite.OAuth2Provider) IntrospectInteractor {
	return IntrospectInteractor{oauth2: oauth2}
}

func (i *IntrospectInteractor) Invoke(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()
	emptySession := models.NewEmptyIdpSession()
	ir, err := i.oauth2.NewIntrospectionRequest(ctx, req, emptySession)
	if err != nil {
		log.Printf("Error occurred in NewIntrospectionRequest: %+v", err)
		i.oauth2.WriteIntrospectionError(ctx, rw, err)
		return err
	}

	i.oauth2.WriteIntrospectionResponse(ctx, rw, ir)
	return nil
}
