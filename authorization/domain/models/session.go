package models

import (
	"time"

	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

func NewSession(clientId, userId string) *openid.DefaultSession {
	header := &jwt.Headers{
		Extra: make(map[string]interface{}),
	}
	header.Add("kid", "ead3f8de")

	claims := &jwt.IDTokenClaims{
		Issuer:      "go-idp",
		Audience:    []string{clientId},
		Subject:     userId,
		IssuedAt:    time.Now(),
		RequestedAt: time.Now(),
		AuthTime:    time.Now(),
	}
	claims.Add("azp", "go-idp")

	return &openid.DefaultSession{
		Claims:  claims,
		Headers: header,
		Subject: userId,
	}
}

func NewEmptySession() *openid.DefaultSession {
	return openid.NewDefaultSession()
}
