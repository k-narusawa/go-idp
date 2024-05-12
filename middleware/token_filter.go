package middleware

import (
	"strings"

	"github.com/k-narusawa/go-idp/resources/application"

	"github.com/labstack/echo/v4"
)

func TokenAuthMiddleware(iu application.IntrospectInteractor) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			splitted := strings.Split(authHeader, " ")
			var bearerToken string

			if len(splitted) == 2 && strings.ToLower(splitted[0]) == "bearer" {
				bearerToken = splitted[1]
			} else {
				return c.JSON(401, map[string]interface{}{
					"error": "invalid_request",
				})
			}

			accessToken, err := iu.Introspect(bearerToken)
			if err != nil {
				return c.JSON(401, map[string]interface{}{
					"error": "invalid_token",
				})
			}

			c.Set("subject", accessToken.Subject)

			return next(c)
		}
	}
}
