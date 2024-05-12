package application

import (
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
)

type RevokeInteractor struct {
	oauth2 fosite.OAuth2Provider
}

func NewRevokeInteractor(oauth2 fosite.OAuth2Provider) RevokeInteractor {
	return RevokeInteractor{oauth2: oauth2}
}

func (r *RevokeInteractor) Invoke(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()
	err := r.oauth2.NewRevocationRequest(ctx, req)

	r.oauth2.WriteRevocationResponse(ctx, rw, err)
	return nil
}
