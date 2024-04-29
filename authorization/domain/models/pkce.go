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

type PKCE struct {
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

func (p *PKCE) SetID(id string) {
	p.ClientID = id
}

func (p *PKCE) GetID() string {
	return p.ClientID
}

func (p *PKCE) GetRequestedAt() time.Time {
	return p.RequestedAt
}

func (p *PKCE) GetClient() fosite.Client {
	return &p.Client
}

func (p *PKCE) GetRequestedScopes() fosite.Arguments {
	return strings.Split(p.Scope, " ")
}

func (p *PKCE) GetRequestedAudience() fosite.Arguments {
	return strings.Split(p.RequestedAudience, " ")
}

func (p *PKCE) SetRequestedScopes(scopes fosite.Arguments) {
	p.Scope = strings.Join(scopes, " ")
}

func (p *PKCE) SetRequestedAudience(audience fosite.Arguments) {
	p.RequestedAudience = strings.Join(audience, " ")
}

func (p *PKCE) AppendRequestedScope(scope string) {
	p.Scope = p.Scope + " " + scope
}

func (p *PKCE) GetGrantedScopes() fosite.Arguments {
	return strings.Split(p.GrantedScope, " ")
}

func (p *PKCE) GetGrantedAudience() fosite.Arguments {
	return strings.Split(p.GrantedAudience, " ")
}

func (p *PKCE) GrantScope(scope string) {
	p.GrantedScope = p.GrantedScope + " " + scope
}

func (p *PKCE) GrantAudience(audience string) {
	p.GrantedAudience = p.GrantedAudience + " " + audience
}

func (p *PKCE) GetSession() fosite.Session {
	var session openid.DefaultSession

	err := json.Unmarshal([]byte(p.SessionData), &session)
	if err != nil {
		log.Printf("Error occurred in GetSession: %+v", err)
		return nil
	}

	return &session
}

func (p *PKCE) SetSession(session fosite.Session) {
	jsonData, err := json.Marshal(session)

	if err != nil {
		return
	}

	p.SessionData = string(jsonData)
}

func (p *PKCE) GetRequestForm() url.Values {
	form, err := url.ParseQuery(p.FormData)
	if err != nil {
		log.Printf("Error occurred in GetRequestForm: %+v", err)
		return nil
	}

	return form
}

func (p *PKCE) Merge(requester fosite.Requester) {
	// Merge implementation goes here
}

func (p *PKCE) Sanitize(allowedParameters []string) fosite.Requester {
	// Sanitize implementation goes here
	return nil
}

func PKCEOf(signature string, requester fosite.Requester) *PKCE {
	jsonData, err := json.Marshal(requester.GetSession())
	if err != nil {
		log.Printf("Error occurred in FromRequester: %+v", err)
		return nil
	}

	return &PKCE{
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

func (p *PKCE) ToRequester() fosite.Requester {
	return &PKCE{
		Signature:         p.Signature,
		ClientID:          p.ClientID,
		Client:            p.Client,
		RequestedAt:       p.RequestedAt,
		Scope:             p.Scope,
		GrantedScope:      p.GrantedScope,
		FormData:          p.FormData,
		Active:            p.Active,
		RequestedAudience: p.RequestedAudience,
		GrantedAudience:   p.GrantedAudience,
		SessionData:       p.SessionData,
	}
}
