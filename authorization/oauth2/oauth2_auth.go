package oauth2

import (
	"idp/authorization/infrastructure"
	"idp/authorization/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthorizationEndpoint(c echo.Context) error {
	rw := c.Response()
	req := c.Request()

	ctx := req.Context()

	// Let's create an AuthorizeRequest object!
	// It will analyze the request and extract important information like scopes, response type and others.
	ar, err := oauth2.NewAuthorizeRequest(ctx, req)
	if err != nil {
		log.Printf("Error occurred in NewAuthorizeRequest: %+v", err)
		oauth2.WriteAuthorizeError(ctx, rw, ar, err)
		return err
	}

	if req.Method == "GET" {
		return c.Render(http.StatusOK, "login.html", nil)
	}

	// req.ParseForm()
	// if req.PostForm.Get("username") != "peter" {
	// 	return c.Render(http.StatusOK, "login.html", nil)
	// }

	for _, scope := range req.PostForm["scopes"] {
		ar.GrantScope(scope)
	}

	un := req.PostForm.Get("username")
	p := req.PostForm.Get("password")

	db := infrastructure.Connect()
	var user models.User
	res := db.Where("username=?", un).First(&user)
	if res.Error != nil {
		log.Printf("Error occurred in GetClient: %+v", res.Error)
		return res.Error
	}

	if err := user.Authenticate(p); err != nil {
		log.Printf("Error occurred in Authenticate: %+v", err)
		return c.Render(http.StatusOK, "login.html", nil)
	}

	mySessionData := newSession(user.UserID)

	// When using the HMACSHA strategy you must use something that implements the HMACSessionContainer.
	// It brings you the power of overriding the default values.
	//
	// mySessionData.HMACSession = &strategy.HMACSession{
	//	AccessTokenExpiry: time.Now().Add(time.Day),
	//	AuthorizeCodeExpiry: time.Now().Add(time.Day),
	// }
	//

	// If you're using the JWT strategy, there's currently no distinction between access token and authorize code claims.
	// Therefore, you both access token and authorize code will have the same "exp" claim. If this is something you
	// need let us know on github.
	//
	// mySessionData.JWTClaims.ExpiresAt = time.Now().Add(time.Day)

	ar.SetResponseTypeHandled("code")
	response, err := oauth2.NewAuthorizeResponse(ctx, ar, mySessionData)

	if err != nil {
		log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
		oauth2.WriteAuthorizeError(ctx, rw, ar, err)
		return err
	}

	oauth2.WriteAuthorizeResponse(ctx, rw, ar, response)

	return nil
}
