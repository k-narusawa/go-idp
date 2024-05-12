package models

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/ory/fosite"
)

type AccessToken struct {
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

func (at *AccessToken) SetID(id string) {
	at.RequestID = id
}

func (at *AccessToken) GetID() string {
	return at.RequestID
}

func (at *AccessToken) GetRequestedAt() time.Time {
	return at.RequestedAt
}

func (at *AccessToken) GetClient() fosite.Client {
	return &at.Client
}

func (at *AccessToken) GetRequestedScopes() fosite.Arguments {
	return strings.Split(at.Scope, " ")
}

func (at *AccessToken) GetRequestedAudience() fosite.Arguments {
	return strings.Split(at.RequestedAudience, " ")
}

func (at *AccessToken) SetRequestedScopes(scopes fosite.Arguments) {
	at.Scope = strings.Join(scopes, " ")
}

func (at *AccessToken) SetRequestedAudience(audience fosite.Arguments) {
	at.RequestedAudience = strings.Join(audience, " ")
}

func (at *AccessToken) AppendRequestedScope(scope string) {
	at.Scope = at.Scope + " " + scope
}

func (at *AccessToken) GetGrantedScopes() fosite.Arguments {
	return strings.Split(at.GrantedScope, " ")
}

func (at *AccessToken) GetGrantedAudience() fosite.Arguments {
	return strings.Split(at.GrantedAudience, " ")
}

func (at *AccessToken) GrantScope(scope string) {
	at.GrantedScope = at.GrantedScope + " " + scope
}

func (at *AccessToken) GrantAudience(audience string) {
	at.GrantedAudience = at.GrantedAudience + " " + audience
}

func (at *AccessToken) IsExpired() bool {
	expiresAt := at.GetSession().GetExpiresAt(fosite.AccessToken)
	return !expiresAt.IsZero() && expiresAt.Before(time.Now())
}

func (at *AccessToken) GetSession() fosite.Session {
	var session fosite.DefaultSession

	if []byte(at.SessionData) == nil || len([]byte(at.SessionData)) == 0 {
		return fosite.NewAuthorizeRequest().Session
	}

	err := json.Unmarshal([]byte(at.SessionData), &session)
	if err != nil {
		log.Printf("Error occurred in GetSession: %+v", err)
		return nil
	}

	return &session
}

func (at *AccessToken) SetSession(session fosite.Session) {
	jsonData, err := json.Marshal(session)

	if err != nil {
		return
	}

	at.SessionData = string(jsonData)
}

func (at *AccessToken) GetRequestForm() url.Values {
	form, err := url.ParseQuery(at.FormData)
	if err != nil {
		log.Printf("Error occurred in GetRequestForm: %+v", err)
		return nil
	}

	return form
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
		RequestID:         requester.GetID(),
		ClientID:          requester.GetClient().GetID(),
		RequestedAt:       requester.GetRequestedAt(),
		Scope:             strings.Join(requester.GetRequestedScopes(), " "),
		GrantedScope:      strings.Join(requester.GetGrantedScopes(), " "),
		FormData:          requester.GetRequestForm().Encode(),
		Active:            true,
		SessionData:       string(jsonData),
		Subject:           requester.GetSession().GetSubject(),
		RequestedAudience: strings.Join(requester.GetRequestedAudience(), " "),
		GrantedAudience:   strings.Join(requester.GetGrantedAudience(), " "),
	}
}

func (at *AccessToken) ToRequester() fosite.Requester {
	return &AccessToken{
		Signature:         at.Signature,
		RequestID:         at.RequestID,
		ClientID:          at.ClientID,
		Client:            at.Client,
		RequestedAt:       at.RequestedAt,
		Scope:             at.Scope,
		GrantedScope:      at.GrantedScope,
		FormData:          at.FormData,
		Active:            at.Active,
		SessionData:       at.SessionData,
		Subject:           at.Subject,
		RequestedAudience: at.RequestedAudience,
		GrantedAudience:   at.GrantedAudience,
	}
}
