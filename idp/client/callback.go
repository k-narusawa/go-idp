package client

import "github.com/labstack/echo/v4"

func CallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")

	if code == "" {
		return c.File("resources/error.html")
	}

	return c.File("resources/success.html")
}
