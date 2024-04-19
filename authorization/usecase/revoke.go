package usecase

import (
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type RevokeUsecase struct {
	oauth2 fosite.OAuth2Provider
}

func NewRevokeUsecase(oauth2 fosite.OAuth2Provider) RevokeUsecase {
	return RevokeUsecase{oauth2: oauth2}
}

func (r *RevokeUsecase) Invoke(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()
	err := r.oauth2.NewRevocationRequest(ctx, req)

	r.oauth2.WriteRevocationResponse(ctx, rw, err)
	return nil
}
