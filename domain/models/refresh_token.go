package models

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/ory/fosite"
)

type RefreshToken struct {
	Signature         string    `gorm:"type:varchar(255);not null;unique" `
	RequestID         string    `gorm:"type:varchar(40);not null"`
	ClientID          string    `gorm:"type:varchar(255);not null"`
	Client            Client    `gorm:"foreignKey:ClientID"`
	RequestedAt       time.Time `gorm:"type:timestamp;not null"`
	Scope             string    `gorm:"type:varchar(255);not null"`
	GrantedScope      string    `gorm:"type:varchar(255);not null"`
	FormData          string    `gorm:"type:text;not null"`
	SessionData       string    `gorm:"type:text;not null"`
	Subject           string    `gorm:"type:varchar(40);not null"`
	Active            bool      `gorm:"type:boolean;not null"`
	RequestedAudience string    `gorm:"type:varchar(255);not null"`
	GrantedAudience   string    `gorm:"type:varchar(255);not null"`
}

func (rt *RefreshToken) SetID(id string) {
	rt.RequestID = id
}

func (rt *RefreshToken) GetID() string {
	return rt.RequestID
}

func (rt *RefreshToken) GetRequestedAt() time.Time {
	return rt.RequestedAt
}

func (rt *RefreshToken) GetClient() fosite.Client {
	return &rt.Client
}

func (rt *RefreshToken) GetRequestedScopes() fosite.Arguments {
	return strings.Split(rt.Scope, " ")
}

func (rt *RefreshToken) GetRequestedAudience() fosite.Arguments {
	return strings.Split(rt.RequestedAudience, " ")
}

func (rt *RefreshToken) SetRequestedScopes(scopes fosite.Arguments) {
	rt.Scope = strings.Join(scopes, " ")
}

func (rt *RefreshToken) SetRequestedAudience(audience fosite.Arguments) {
	rt.RequestedAudience = strings.Join(audience, " ")
}

func (rt *RefreshToken) AppendRequestedScope(scope string) {
	rt.Scope = rt.Scope + " " + scope
}

func (rt *RefreshToken) GetGrantedScopes() fosite.Arguments {
	return strings.Split(rt.GrantedScope, " ")
}

func (rt *RefreshToken) GetGrantedAudience() fosite.Arguments {
	return strings.Split(rt.GrantedAudience, " ")
}

func (rt *RefreshToken) GrantScope(scope string) {
	rt.GrantedScope = rt.GrantedScope + " " + scope
}

func (rt *RefreshToken) GrantAudience(audience string) {
	rt.GrantedAudience = rt.GrantedAudience + " " + audience
}

func (rt *RefreshToken) GetSession() fosite.Session {
	var session fosite.DefaultSession

	err := json.Unmarshal([]byte(rt.SessionData), &session)
	if err != nil {
		log.Printf("Error occurred in GetSession: %+v", err)
		return nil
	}

	return &session
}

func (rt *RefreshToken) SetSession(session fosite.Session) {
	jsonData, err := json.Marshal(session)

	if err != nil {
		return
	}

	rt.SessionData = string(jsonData)
}

func (rt *RefreshToken) GetRequestForm() url.Values {
	form, err := url.ParseQuery(rt.FormData)
	if err != nil {
		log.Printf("Error occurred in GetRequestForm: %+v", err)
		return nil
	}

	return form
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
		RequestID:         requester.GetID(),
		ClientID:          requester.GetClient().GetID(),
		RequestedAt:       requester.GetRequestedAt(),
		Scope:             strings.Join(requester.GetRequestedScopes(), " "),
		GrantedScope:      strings.Join(requester.GetGrantedScopes(), " "),
		FormData:          requester.GetRequestForm().Encode(),
		SessionData:       string(jsonData),
		Subject:           requester.GetSession().GetSubject(),
		Active:            true,
		RequestedAudience: strings.Join(requester.GetRequestedAudience(), " "),
		GrantedAudience:   strings.Join(requester.GetGrantedAudience(), " "),
	}
}

func (rt *RefreshToken) ToRequester() fosite.Requester {
	return &RefreshToken{
		Signature:         rt.Signature,
		RequestID:         rt.RequestID,
		ClientID:          rt.ClientID,
		Client:            rt.Client,
		RequestedAt:       rt.RequestedAt,
		Scope:             rt.Scope,
		GrantedScope:      rt.GrantedScope,
		FormData:          rt.FormData,
		SessionData:       rt.SessionData,
		Subject:           rt.Subject,
		Active:            rt.Active,
		RequestedAudience: rt.RequestedAudience,
		GrantedAudience:   rt.GrantedAudience,
	}
}
