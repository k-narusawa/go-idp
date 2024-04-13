package models

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"gorm.io/gorm"
)

type IDSession struct {
	gorm.Model
	Signature         string    `gorm:"type:varchar(255);not null;unique" `
	ClientID          string    `gorm:"type:varchar(255);not null"`
	Client            Client    `gorm:"foreignKey:ClientID"`
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

func (is *IDSession) SetID(id string) {
	is.ClientID = id
}

func (is *IDSession) GetID() string {
	return is.ClientID
}

func (is *IDSession) GetRequestedAt() time.Time {
	return is.RequestedAt
}

func (is *IDSession) GetClient() fosite.Client {
	return &is.Client
}

func (is *IDSession) GetRequestedScopes() fosite.Arguments {
	return strings.Split(is.Scope, " ")
}

func (is *IDSession) GetRequestedAudience() fosite.Arguments {
	return strings.Split(is.RequestedAudience, " ")
}

func (is *IDSession) SetRequestedScopes(scopes fosite.Arguments) {
	is.Scope = strings.Join(scopes, " ")
}

func (is *IDSession) SetRequestedAudience(audience fosite.Arguments) {
	is.RequestedAudience = strings.Join(audience, " ")
}

func (is *IDSession) AppendRequestedScope(scope string) {
	is.Scope = is.Scope + " " + scope
}

func (is *IDSession) GetGrantedScopes() fosite.Arguments {
	return strings.Split(is.GrantedScope, " ")
}

func (is *IDSession) GetGrantedAudience() fosite.Arguments {
	return strings.Split(is.GrantedAudience, " ")
}

func (is *IDSession) GrantScope(scope string) {
	is.GrantedScope = is.GrantedScope + " " + scope
}

func (is *IDSession) GrantAudience(audience string) {
	is.GrantedAudience = is.GrantedAudience + " " + audience
}

func (is *IDSession) GetSession() fosite.Session {
	var session openid.DefaultSession

	err := json.Unmarshal([]byte(is.SessionData), &session)
	if err != nil {
		log.Printf("Error occurred in GetSession: %+v", err)
		return nil
	}

	return &session
}

func (is *IDSession) SetSession(session fosite.Session) {
	jsonData, err := json.Marshal(session)

	if err != nil {
		return
	}

	is.SessionData = string(jsonData)
}

func (is *IDSession) GetRequestForm() url.Values {
	return url.Values{}
}

func (is *IDSession) Merge(requester fosite.Requester) {
	// Merge implementation goes here
}

func (is *IDSession) Sanitize(allowedParameters []string) fosite.Requester {
	// Sanitize implementation goes here
	return nil
}

func IDSessionOf(signature string, requester fosite.Requester) *IDSession {
	jsonData, err := json.Marshal(requester.GetSession())
	if err != nil {
		log.Printf("Error occurred in FromRequester: %+v", err)
		return nil
	}

	return &IDSession{
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

func (is *IDSession) ToRequester() fosite.Requester {
	return &IDSession{
		Signature:         is.Signature,
		ClientID:          is.ClientID,
		Client:            is.Client,
		RequestedAt:       is.RequestedAt,
		Scope:             is.Scope,
		GrantedScope:      is.GrantedScope,
		FormData:          is.FormData,
		Active:            is.Active,
		RequestedAudience: is.RequestedAudience,
		GrantedAudience:   is.GrantedAudience,
		SessionData:       is.SessionData,
	}
}
