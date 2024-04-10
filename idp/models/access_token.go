package models

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/ory/fosite"
	"gorm.io/gorm"
)

type AccessToken struct {
	gorm.Model
	Signature         string    `gorm:"type:varchar(255);not null;unique" `
	clientID          string    `gorm:"type:varchar(255);not null"`
	requestedAt       time.Time `gorm:"type:timestamp;not null"`
	scope             string    `gorm:"type:varchar(255);not null"`
	GrantedScope      string    `gorm:"type:varchar(255);not null"`
	FormData          string    `gorm:"type:text;not null"`
	SessionData       string    `gorm:"type:varchar(255);not null"`
	Subject           string    `gorm:"type:text;not null"`
	Active            bool      `gorm:"type:boolean;not null"`
	requestedAudience string    `gorm:"type:varchar(255);not null"`
	grantedAudience   string    `gorm:"type:varchar(255);not null"`
}

func (a *AccessToken) SetID(id string) {
	a.clientID = id
}

func (a *AccessToken) GetID() string {
	return a.clientID
}

func (a *AccessToken) GetRequestedAt() time.Time {
	return a.requestedAt
}

func (a *AccessToken) GetClient() fosite.Client {
	return fosite.Client(&Client{ID: a.clientID})
}

func (a *AccessToken) GetRequestedScopes() fosite.Arguments {
	return strings.Split(a.scope, " ")
}

func (a *AccessToken) GetRequestedAudience() fosite.Arguments {
	return strings.Split(a.requestedAudience, " ")
}

func (a *AccessToken) SetRequestedScopes(scopes fosite.Arguments) {
	a.scope = strings.Join(scopes, " ")
}

func (a *AccessToken) SetRequestedAudience(audience fosite.Arguments) {
	a.requestedAudience = strings.Join(audience, " ")
}

func (a *AccessToken) AppendRequestedScope(scope string) {
	a.scope = a.scope + " " + scope
}

func (a *AccessToken) GetGrantedScopes() fosite.Arguments {
	return strings.Split(a.GrantedScope, " ")
}

func (a *AccessToken) GetGrantedAudience() fosite.Arguments {
	return strings.Split(a.grantedAudience, " ")
}

func (a *AccessToken) GrantScope(scope string) {
	a.GrantedScope = a.GrantedScope + " " + scope
}

func (a *AccessToken) GrantAudience(audience string) {
	a.grantedAudience = a.grantedAudience + " " + audience
}

func (a *AccessToken) GetSession() fosite.Session {
	var session fosite.DefaultSession

	err := json.Unmarshal([]byte(a.SessionData), &session)
	if err != nil {
		log.Printf("Error occurred in GetSession: %+v", err)
		return nil
	}

	log.Printf("session: %+v", session)

	return &session
}

func (a *AccessToken) SetSession(session fosite.Session) {
	jsonData, err := json.Marshal(session)

	if err != nil {
		return
	}

	a.SessionData = string(jsonData)
}

func (a *AccessToken) GetRequestForm() url.Values {
	return url.Values{}
}

func (a *AccessToken) Merge(requester fosite.Requester) {
	// Merge implementation goes here
}

func (a *AccessToken) Sanitize(allowedParameters []string) fosite.Requester {
	// Sanitize implementation goes here
	return nil
}

func AccessTokenOf(signature string, requester fosite.Requester) *AccessToken {
	jsonData, err := json.Marshal(requester.GetSession())
	if err != nil {
		log.Printf("Error occurred in FromRequester: %+v", err)
		return nil
	}

	return &AccessToken{
		Signature:         signature,
		clientID:          requester.GetClient().GetID(),
		requestedAt:       requester.GetRequestedAt(),
		scope:             strings.Join(requester.GetRequestedScopes(), " "),
		GrantedScope:      strings.Join(requester.GetGrantedScopes(), " "),
		FormData:          requester.GetRequestForm().Encode(),
		Active:            true,
		requestedAudience: strings.Join(requester.GetRequestedAudience(), " "),
		grantedAudience:   strings.Join(requester.GetGrantedAudience(), " "),
		SessionData:       string(jsonData),
	}
}

func (a *AccessToken) ToRequester() fosite.Requester {
	return &AccessToken{
		Signature:         a.Signature,
		clientID:          a.clientID,
		requestedAt:       a.requestedAt,
		scope:             a.scope,
		GrantedScope:      a.GrantedScope,
		FormData:          a.FormData,
		Active:            a.Active,
		requestedAudience: a.requestedAudience,
		grantedAudience:   a.grantedAudience,
		SessionData:       a.SessionData,
	}
}
