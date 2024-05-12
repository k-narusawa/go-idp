package application

import (
	"github.com/labstack/echo/v4"
)

type UserinfoInteractor struct {
}

func NewUserinfoInteractor() UserinfoInteractor {
	return UserinfoInteractor{}
}

func (ui *UserinfoInteractor) GetUserinfo(c echo.Context) error {
	sub := c.Get(("subject")).(string)

	return c.JSON(200, map[string]interface{}{
		"sub": sub,
	})
}
