package models

import (
	"time"

	"github.com/mohae/deepcopy"
	"github.com/ory/fosite"
	"github.com/ory/fosite/token/jwt"
)

type Session interface {
	IDTokenClaims() *jwt.IDTokenClaims
	IDTokenHeaders() *jwt.Headers
	fosite.Session
}

type IdpSession struct {
	SessionID string                         `json:"session_id"`
	Claims    *jwt.IDTokenClaims             `json:"id_token_claims"`
	Headers   *jwt.Headers                   `json:"headers"`
	ExpiresAt map[fosite.TokenType]time.Time `json:"expires_at"`
	Username  string                         `json:"username"`
	Subject   string                         `json:"subject"`
}

func NewIdpSession(clientId, userId string) *IdpSession {
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

	return &IdpSession{
		Claims:  claims,
		Headers: header,
		Subject: userId,
	}
}

func NewEmptyIdpSession() *IdpSession {
	return &IdpSession{
		Claims:  &jwt.IDTokenClaims{},
		Headers: &jwt.Headers{},
	}
}

func (s *IdpSession) SetSessionID(id string) {
	s.SessionID = id
}

func (s *IdpSession) SetExpiresAt(key fosite.TokenType, exp time.Time) {
	if s.ExpiresAt == nil {
		s.ExpiresAt = make(map[fosite.TokenType]time.Time)
	}
	s.ExpiresAt[key] = exp
}

func (s *IdpSession) GetExpiresAt(key fosite.TokenType) time.Time {
	if s.ExpiresAt == nil {
		s.ExpiresAt = make(map[fosite.TokenType]time.Time)
	}

	if _, ok := s.ExpiresAt[key]; !ok {
		return time.Time{}
	}
	return s.ExpiresAt[key]
}

func (s *IdpSession) GetUsername() string {
	if s == nil {
		return ""
	}
	return s.Username
}

func (s *IdpSession) SetSubject(subject string) {
	s.Subject = subject
}

func (s *IdpSession) GetSubject() string {
	if s == nil {
		return ""
	}

	return s.Subject
}

func (s *IdpSession) Clone() fosite.Session {
	if s == nil {
		return nil
	}

	return deepcopy.Copy(s).(fosite.Session)
}

func (s *IdpSession) IDTokenHeaders() *jwt.Headers {
	if s.Headers == nil {
		s.Headers = &jwt.Headers{}
	}
	return s.Headers
}

func (s *IdpSession) IDTokenClaims() *jwt.IDTokenClaims {
	if s.Claims == nil {
		s.Claims = &jwt.IDTokenClaims{}
	}
	return s.Claims
}
