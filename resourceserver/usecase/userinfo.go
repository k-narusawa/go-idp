package usecase

import (
	"idp/resourceserver/domain/models"

	"github.com/labstack/echo/v4"
)

type UserinfoUsecase struct {
}

func NewAuthorization() UserinfoUsecase {
	return UserinfoUsecase{}
}

func (ui *UserinfoUsecase) GetUserinfo(c echo.Context) error {
	ir := c.Get(("ir")).(models.IntrospectResponse)

	return c.JSON(200, map[string]interface{}{
		"sub": ir.Sub,
	})
}
