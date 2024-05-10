package middleware

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	model "github.com/k-narusawa/go-idp/resources/domain/models"

	"github.com/labstack/echo/v4"
)

func TokenAuthMiddleware() echo.MiddlewareFunc {
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

			values := url.Values{}
			values.Set("token", bearerToken)
			te := "http://localhost:3846/oauth2/introspect"

			req, _ := http.NewRequest(http.MethodPost, te, strings.NewReader(values.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.SetBasicAuth("go-idp", "%7EEp52-Sp%25iQtcEHpSLQ5%2CLT-%2C9%2AHMNfg%2C7WP")

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				return err
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return c.JSON(401, map[string]interface{}{
					"error": "invalid_token",
				})
			}

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			var ir model.IntrospectResponse
			err = json.Unmarshal(b, &ir)
			if err != nil {
				log.Fatal(err)
			}

			c.Set("ir", ir)

			return next(c)
		}
	}
}
