package models

import (
	"time"

	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

func NewSession(userId string) *openid.DefaultSession {
	header := &jwt.Headers{
		Extra: make(map[string]interface{}),
	}
	header.Add("kid", "go-idp:123")

	claims := &jwt.IDTokenClaims{
		Issuer:      "go-idp",
		Audience:    []string{"my-client"},
		Subject:     userId,
		IssuedAt:    time.Now(),
		RequestedAt: time.Now(),
		AuthTime:    time.Now(),
	}
	claims.Add("azp", "my-client")

	return &openid.DefaultSession{
		Claims:  claims,
		Headers: header,
		Subject: userId,
	}
}
