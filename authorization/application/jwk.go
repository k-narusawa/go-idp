package application

import (
	"encoding/json"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type JWKInteractor struct{}

func NewJWKInteractor() JWKInteractor {
	return JWKInteractor{}
}

func (j *JWKInteractor) Invoke(c echo.Context) error {
	data, _ := os.ReadFile("cert/public.pem")
	keyset, _ := jwk.ParseKey(data, jwk.WithPEM(true))

	keyset.Set(jwk.KeyIDKey, "ead3f8de")

	jwk := map[string]interface{}{
		"keys": []interface{}{keyset},
	}
	buf, _ := json.MarshalIndent(jwk, "", "  ")

	return c.JSONBlob(200, buf)
}
