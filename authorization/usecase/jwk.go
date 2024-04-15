package usecase

import (
	"encoding/json"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type JWKUsecase struct{}

func NewJWKUsecase() JWKUsecase {
	return JWKUsecase{}
}

func (j *JWKUsecase) Invoke(c echo.Context) error {
	data, _ := os.ReadFile("cert/public-key.pem")
	keyset, _ := jwk.ParseKey(data, jwk.WithPEM(true))

	keyset.Set(jwk.KeyIDKey, "go-idp:123")

	jwk := map[string]interface{}{
		"keys": []interface{}{keyset},
	}
	buf, _ := json.MarshalIndent(jwk, "", "  ")

	return c.JSONBlob(200, buf)
}
