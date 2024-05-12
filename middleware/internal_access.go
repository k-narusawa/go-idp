package middleware

import (
	"net/http"

	"github.com/k-narusawa/go-idp/util"
	"github.com/labstack/echo/v4"
)

func InternalAccess() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			networkUtil := util.NewNetworkUtil()

			isPrivateAddress := networkUtil.IsPrivateAddress(c.RealIP())

			if !isPrivateAddress {
				return c.String(http.StatusForbidden, "Forbidden")
			}

			return next(c)
		}
	}
}
