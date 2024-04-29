package usecase

import (
	"log"

	"github.com/k-narusawa/go-idp/authorization/domain/models"

	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type IntrospectUsecase struct {
	oauth2 fosite.OAuth2Provider
}

func NewIntrospectUsecase(oauth2 fosite.OAuth2Provider) IntrospectUsecase {
	return IntrospectUsecase{oauth2: oauth2}
}

func (i *IntrospectUsecase) Invoke(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()
	mySessionData := models.NewSession("")
	ir, err := i.oauth2.NewIntrospectionRequest(ctx, req, mySessionData)
	if err != nil {
		log.Printf("Error occurred in NewIntrospectionRequest: %+v", err)
		i.oauth2.WriteIntrospectionError(ctx, rw, err)
		return err
	}

	i.oauth2.WriteIntrospectionResponse(ctx, rw, ir)
	return nil
}
