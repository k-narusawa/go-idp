package client

import (
	"encoding/json"
	"idp/models"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

func CallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")

	if code == "" {
		return c.Render(http.StatusOK, "error.html", nil)
	}

	values := url.Values{}
	values.Set("code", code)
	values.Add("redirect_uri", "http://localhost:3846/callback")
	values.Add("grant_type", "authorization_code")
	te := "http://localhost:3846/oauth2/token"

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
		return c.Render(http.StatusOK, "error.html", nil)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tr models.TokenResponse
	err = json.Unmarshal(b, &tr)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, tr)
}
