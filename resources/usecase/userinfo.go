package usecase

import (
	"github.com/labstack/echo/v4"
)

type UserinfoUsecase struct {
}

func NewAuthorization() UserinfoUsecase {
	return UserinfoUsecase{}
}

func (ui *UserinfoUsecase) GetUserinfo(c echo.Context) error {
	sub := c.Get(("subject")).(string)

	return c.JSON(200, map[string]interface{}{
		"sub": sub,
	})
}
