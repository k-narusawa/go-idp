package oauth2

import (
	"log"

	"github.com/labstack/echo/v4"
)

func IntrospectionEndpoint(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()
	mySessionData := newSession("")
	ir, err := oauth2.NewIntrospectionRequest(ctx, req, mySessionData)
	if err != nil {
		log.Printf("Error occurred in NewIntrospectionRequest: %+v", err)
		oauth2.WriteIntrospectionError(ctx, rw, err)
		return err
	}

	oauth2.WriteIntrospectionResponse(ctx, rw, ir)
	return nil
}
