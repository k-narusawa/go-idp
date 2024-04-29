package models

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
)

type AuthorizationCode struct {
	Signature         string    `gorm:"type:varchar(255);not null;unique" `
	ClientID          string    `gorm:"type:varchar(255);not null"`
	Client            Client    `gorm:"foreignKey:ClientID;"`
	RequestedAt       time.Time `gorm:"type:timestamp;not null"`
	Scope             string    `gorm:"type:varchar(255);not null"`
	GrantedScope      string    `gorm:"type:varchar(255);not null"`
	FormData          string    `gorm:"type:text;not null"`
	SessionData       string    `gorm:"type:text;not null"`
	Subject           string    `gorm:"type:text;not null"`
	Active            bool      `gorm:"type:boolean;not null"`
	RequestedAudience string    `gorm:"type:varchar(255);not null"`
	GrantedAudience   string    `gorm:"type:varchar(255);not null"`
}

func (ac *AuthorizationCode) SetID(id string) {
	ac.ClientID = id
}

func (ac *AuthorizationCode) GetID() string {
	return ac.ClientID
}

func (ac *AuthorizationCode) GetClient() fosite.Client {
	return &ac.Client
}

func (ac *AuthorizationCode) SetClient(c Client) {
	ac.Client = c
}

func (ac *AuthorizationCode) GetRequestedAt() (requestedAt time.Time) {
	return requestedAt
}

func (ac *AuthorizationCode) GetRequestedScopes() fosite.Arguments {
	return strings.Split(ac.Scope, " ")
}

func (ac *AuthorizationCode) GetRequestedAudience() fosite.Arguments {
	return strings.Split(ac.RequestedAudience, " ")
}

func (ac *AuthorizationCode) SetRequestedScopes(scopes fosite.Arguments) {
	ac.Scope = strings.Join(scopes, " ")
}

func (ac *AuthorizationCode) SetRequestedAudience(audience fosite.Arguments) {
	ac.RequestedAudience = strings.Join(audience, " ")
}

func (ac *AuthorizationCode) AppendRequestedScope(scope string) {
	ac.Scope = ac.Scope + " " + scope
}

func (ac *AuthorizationCode) GetGrantedScopes() fosite.Arguments {
	return strings.Split(ac.GrantedScope, " ")
}

func (ac *AuthorizationCode) GetGrantedAudience() fosite.Arguments {
	return strings.Split(ac.GrantedAudience, " ")
}

func (ac *AuthorizationCode) GrantScope(scope string) {
	ac.GrantedScope = ac.GrantedScope + " " + scope
}

func (ac *AuthorizationCode) GrantAudience(audience string) {
	ac.GrantedAudience = ac.GrantedAudience + " " + audience
}

func (ac *AuthorizationCode) GetSession() fosite.Session {
	var session openid.DefaultSession

	err := json.Unmarshal([]byte(ac.SessionData), &session)
	if err != nil {
		log.Printf("Error occurred in GetSession: %+v", err)
		return nil
	}

	return &session
}

func (ac *AuthorizationCode) SetSession(session fosite.Session) {
	jsonData, err := json.Marshal(session)

	if err != nil {
		return
	}

	ac.SessionData = string(jsonData)
}

func (ac *AuthorizationCode) GetRequestForm() url.Values {
	form, err := url.ParseQuery(ac.FormData)
	if err != nil {
		log.Printf("Error occurred in GetRequestForm: %+v", err)
		return nil
	}

	return form
}

func (ac *AuthorizationCode) Merge(requester fosite.Requester) {
	// Merge implementation goes here
}

func (ac *AuthorizationCode) Sanitize(allowedParameters []string) fosite.Requester {
	// Sanitize implementation goes here
	return nil
}

func AuthorizationCodeOf(signature string, requester fosite.Requester) *AuthorizationCode {
	jsonData, err := json.Marshal(requester.GetSession())
	if err != nil {
		log.Printf("Error occurred in FromRequester: %+v", err)
		return nil
	}

	return &AuthorizationCode{
		Signature:         signature,
		ClientID:          requester.GetClient().GetID(),
		RequestedAt:       requester.GetRequestedAt(),
		Scope:             strings.Join(requester.GetRequestedScopes(), " "),
		GrantedScope:      strings.Join(requester.GetGrantedScopes(), " "),
		FormData:          requester.GetRequestForm().Encode(),
		Active:            true,
		RequestedAudience: strings.Join(requester.GetRequestedAudience(), " "),
		GrantedAudience:   strings.Join(requester.GetGrantedAudience(), " "),
		SessionData:       string(jsonData),
	}
}

func (ac *AuthorizationCode) ToRequester() fosite.Requester {
	return &AuthorizationCode{
		Signature:         ac.Signature,
		ClientID:          ac.ClientID,
		Client:            ac.Client,
		RequestedAt:       ac.RequestedAt,
		Scope:             ac.Scope,
		GrantedScope:      ac.GrantedScope,
		FormData:          ac.FormData,
		Active:            ac.Active,
		RequestedAudience: ac.RequestedAudience,
		GrantedAudience:   ac.GrantedAudience,
		SessionData:       ac.SessionData,
	}
}
