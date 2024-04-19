package usecase

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserinfoUsecase struct {
}

func NewAuthorization() UserinfoUsecase {
	return UserinfoUsecase{}
}

func (ui *UserinfoUsecase) GetUserinfo(c echo.Context) error {
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
	req.SetBasicAuth("my-client", "foobar")

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

	var ir IntrospectResponse
	err = json.Unmarshal(b, &ir)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(200, map[string]interface{}{
		"sub": ir.Sub,
	})
}

type IntrospectResponse struct {
	Active   bool     `json:"active"`
	Aud      []string `json:"aud"`
	ClientId string   `json:"client_id"`
	Exp      int      `json:"exp"`
	Iat      int      `json:"iat"`
	Scope    string   `json:"scope"`
	Sub      string   `json:"sub"`
}
