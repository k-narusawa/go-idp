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

type RefreshToken struct {
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

func (rt *RefreshToken) SetID(id string) {
	rt.clientID = id
}

func (rt *RefreshToken) GetID() string {
	return rt.clientID
}

func (rt *RefreshToken) GetRequestedAt() time.Time {
	return rt.requestedAt
}

func (rt *RefreshToken) GetClient() fosite.Client {
	return fosite.Client(&Client{ID: rt.clientID})
}

func (rt *RefreshToken) GetRequestedScopes() fosite.Arguments {
	return strings.Split(rt.scope, " ")
}

func (rt *RefreshToken) GetRequestedAudience() fosite.Arguments {
	return strings.Split(rt.requestedAudience, " ")
}

func (rt *RefreshToken) SetRequestedScopes(scopes fosite.Arguments) {
	rt.scope = strings.Join(scopes, " ")
}

func (rt *RefreshToken) SetRequestedAudience(audience fosite.Arguments) {
	rt.requestedAudience = strings.Join(audience, " ")
}

func (rt *RefreshToken) AppendRequestedScope(scope string) {
	rt.scope = rt.scope + " " + scope
}

func (rt *RefreshToken) GetGrantedScopes() fosite.Arguments {
	return strings.Split(rt.grantedScope, " ")
}

func (rt *RefreshToken) GetGrantedAudience() fosite.Arguments {
	return strings.Split(rt.grantedAudience, " ")
}

func (rt *RefreshToken) GrantScope(scope string) {
	rt.grantedScope = rt.grantedScope + " " + scope
}

func (rt *RefreshToken) GrantAudience(audience string) {
	rt.grantedAudience = rt.grantedAudience + " " + audience
}

func (rt *RefreshToken) GetSession() fosite.Session {
	var session fosite.DefaultSession

	err := json.Unmarshal([]byte(rt.sessionData), &session)
	if err != nil {
		log.Printf("Error occurred in GetSession: %+v", err)
		return nil
	}

	log.Printf("session: %+v", session)

	return &session
}

func (rt *RefreshToken) SetSession(session fosite.Session) {
	jsonData, err := json.Marshal(session)

	if err != nil {
		return
	}

	rt.sessionData = string(jsonData)
}

func (rt *RefreshToken) GetRequestForm() url.Values {
	return url.Values{}
}

func (rt *RefreshToken) Merge(requester fosite.Requester) {
	// Merge implementation goes here
}

func (rt *RefreshToken) Sanitize(allowedParameters []string) fosite.Requester {
	// Sanitize implementation goes here
	return nil
}

func RefreshTokenOf(signature string, requester fosite.Requester) *RefreshToken {
	jsonData, err := json.Marshal(requester.GetSession())
	if err != nil {
		log.Printf("Error occurred in FromRequester: %+v", err)
		return nil
	}

	return &RefreshToken{
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

func (rt *RefreshToken) ToRequester() fosite.Requester {
	return &RefreshToken{
		Signature:         rt.Signature,
		clientID:          rt.clientID,
		requestedAt:       rt.requestedAt,
		scope:             rt.scope,
		grantedScope:      rt.grantedScope,
		formData:          rt.formData,
		active:            rt.active,
		requestedAudience: rt.requestedAudience,
		grantedAudience:   rt.grantedAudience,
		sessionData:       rt.sessionData,
	}
}
