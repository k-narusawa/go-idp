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
	data, _ := os.ReadFile("cert/public.pem")
	keyset, _ := jwk.ParseKey(data, jwk.WithPEM(true))

	keyset.Set(jwk.KeyIDKey, "go-idp:123")

	jwk := map[string]interface{}{
		"keys": []interface{}{keyset},
	}
	buf, _ := json.MarshalIndent(jwk, "", "  ")

	return c.JSONBlob(200, buf)
	// str := `{
	// 	"keys": [
	// 		{
	// 			"e": "AQAX",
	// 			"kid": "ftGzDRLZemP6pqS14vrIvEcrUiAW7MBXdSwtAAJok_k",
	// 			"kty": "RSA",
	// 			"n": "0LcF_hKvj4JPVbyvS1pZMddUcjk6dLDLvIWC-hQ4ZAwRG6lD_L0Y0nUiXST3FsXlBn4ivGb8y4dH2puXDS-CvA5e4zaZGk7P6ypNEEBGfcoN5BxhAp9zU7fPcLWU9eSkYeKCYJwPx5Wk8ohlRxAqzJFt3f41BwbMHBgzDOCyyIeE47W6cynPrgf8MmU7vJIT5mNHnBVXC_ktSjVL86nQJ19Sx3h9Au2CNw1iys2yD9u-wx987da25O2DzaegXZ4bj5IrfgdIdz-TlDYvqalx5KheX5ZimbcDAYd0AjaR-h8p153oJBsMdDkwoMVY55Qs5Nsw_5PlQ80Kr6TQJPwm1Q"
	// 		}
	// 	]
	// }
	// `

	// // 文字列をJSONに変換
	// var jwkMap map[string]interface{}
	// if err := json.Unmarshal([]byte(str), &jwkMap); err != nil {
	// 	log.Fatal(err)
	// }

	// return c.JSON(200, jwkMap)
}
