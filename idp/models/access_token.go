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
	grantedScope      string    `gorm:"type:varchar(255);not null"`
	formData          string    `gorm:"type:text;not null"`
	sessionData       string    `gorm:"type:varchar(255);not null"`
	Subject           string    `gorm:"type:text;not null"`
	active            bool      `gorm:"type:boolean;not null"`
	requestedAudience string    `gorm:"type:varchar(255);not null"`
	grantedAudience   string    `gorm:"type:varchar(255);not null"`
}

func (at *AccessToken) SetID(id string) {
	at.clientID = id
}

func (at *AccessToken) GetID() string {
	return at.clientID
}

func (at *AccessToken) GetRequestedAt() time.Time {
	return at.requestedAt
}

func (at *AccessToken) GetClient() fosite.Client {
	return fosite.Client(&Client{ID: at.clientID})
}

func (at *AccessToken) GetRequestedScopes() fosite.Arguments {
	return strings.Split(at.scope, " ")
}

func (at *AccessToken) GetRequestedAudience() fosite.Arguments {
	return strings.Split(at.requestedAudience, " ")
}

func (at *AccessToken) SetRequestedScopes(scopes fosite.Arguments) {
	at.scope = strings.Join(scopes, " ")
}

func (at *AccessToken) SetRequestedAudience(audience fosite.Arguments) {
	at.requestedAudience = strings.Join(audience, " ")
}

func (at *AccessToken) AppendRequestedScope(scope string) {
	at.scope = at.scope + " " + scope
}

func (at *AccessToken) GetGrantedScopes() fosite.Arguments {
	return strings.Split(at.grantedScope, " ")
}

func (at *AccessToken) GetGrantedAudience() fosite.Arguments {
	return strings.Split(at.grantedAudience, " ")
}

func (at *AccessToken) GrantScope(scope string) {
	at.grantedScope = at.grantedScope + " " + scope
}

func (at *AccessToken) GrantAudience(audience string) {
	at.grantedAudience = at.grantedAudience + " " + audience
}

func (at *AccessToken) GetSession() fosite.Session {
	var session fosite.DefaultSession

	err := json.Unmarshal([]byte(at.sessionData), &session)
	if err != nil {
		log.Printf("Error occurred in GetSession: %+v", err)
		return nil
	}

	log.Printf("session: %+v", session)

	return &session
}

func (at *AccessToken) SetSession(session fosite.Session) {
	jsonData, err := json.Marshal(session)

	if err != nil {
		return
	}

	at.sessionData = string(jsonData)
}

func (at *AccessToken) GetRequestForm() url.Values {
	return url.Values{}
}

func (at *AccessToken) Merge(requester fosite.Requester) {
	// Merge implementation goes here
}

func (at *AccessToken) Sanitize(allowedParameters []string) fosite.Requester {
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
		grantedScope:      strings.Join(requester.GetGrantedScopes(), " "),
		formData:          requester.GetRequestForm().Encode(),
		active:            true,
		requestedAudience: strings.Join(requester.GetRequestedAudience(), " "),
		grantedAudience:   strings.Join(requester.GetGrantedAudience(), " "),
		sessionData:       string(jsonData),
	}
}

func (at *AccessToken) ToRequester() fosite.Requester {
	return &AccessToken{
		Signature:         at.Signature,
		clientID:          at.clientID,
		requestedAt:       at.requestedAt,
		scope:             at.scope,
		grantedScope:      at.grantedScope,
		formData:          at.formData,
		active:            at.active,
		requestedAudience: at.requestedAudience,
		grantedAudience:   at.grantedAudience,
		sessionData:       at.sessionData,
	}
}
