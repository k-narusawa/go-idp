package usecase

import (
	"idp/resourceserver/domain/model"

	"github.com/labstack/echo/v4"
)

type UserinfoUsecase struct {
}

func NewAuthorization() UserinfoUsecase {
	return UserinfoUsecase{}
}

func (ui *UserinfoUsecase) GetUserinfo(c echo.Context) error {
	ir := c.Get(("ir")).(model.IntrospectResponse)

	return c.JSON(200, map[string]interface{}{
		"sub": ir.Sub,
	})
}
