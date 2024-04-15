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

	return &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Issuer:      "go-idp",
			Subject:     userId,
			IssuedAt:    time.Now(),
			RequestedAt: time.Now(),
			AuthTime:    time.Now(),
		},
		Headers: header,
		Subject: userId,
	}
}
