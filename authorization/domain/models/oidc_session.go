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

type OidcSession struct {
	Signature         string    `gorm:"type:varchar(255);not null;unique" `
	RequestID         string    `gorm:"type:varchar(40);not null"`
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

func (is *OidcSession) SetID(id string) {
	is.RequestID = id
}

func (is *OidcSession) GetID() string {
	return is.RequestID
}

func (is *OidcSession) GetRequestedAt() time.Time {
	return is.RequestedAt
}

func (is *OidcSession) GetClient() fosite.Client {
	return &is.Client
}

func (is *OidcSession) GetRequestedScopes() fosite.Arguments {
	return strings.Split(is.Scope, " ")
}

func (is *OidcSession) GetRequestedAudience() fosite.Arguments {
	return strings.Split(is.RequestedAudience, " ")
}

func (is *OidcSession) SetRequestedScopes(scopes fosite.Arguments) {
	is.Scope = strings.Join(scopes, " ")
}

func (is *OidcSession) SetRequestedAudience(audience fosite.Arguments) {
	is.RequestedAudience = strings.Join(audience, " ")
}

func (is *OidcSession) AppendRequestedScope(scope string) {
	is.Scope = is.Scope + " " + scope
}

func (is *OidcSession) GetGrantedScopes() fosite.Arguments {
	return strings.Split(is.GrantedScope, " ")
}

func (is *OidcSession) GetGrantedAudience() fosite.Arguments {
	return strings.Split(is.GrantedAudience, " ")
}

func (is *OidcSession) GrantScope(scope string) {
	is.GrantedScope = is.GrantedScope + " " + scope
}

func (is *OidcSession) GrantAudience(audience string) {
	is.GrantedAudience = is.GrantedAudience + " " + audience
}

func (is *OidcSession) GetSession() fosite.Session {
	var session openid.DefaultSession

	err := json.Unmarshal([]byte(is.SessionData), &session)
	if err != nil {
		log.Printf("Error occurred in GetSession: %+v", err)
		return nil
	}

	return &session
}

func (is *OidcSession) SetSession(session fosite.Session) {
	jsonData, err := json.Marshal(session)

	if err != nil {
		return
	}

	is.SessionData = string(jsonData)
}

func (is *OidcSession) GetRequestForm() url.Values {
	form, err := url.ParseQuery(is.FormData)
	if err != nil {
		log.Printf("Error occurred in GetRequestForm: %+v", err)
		return nil
	}

	return form
}

func (is *OidcSession) Merge(requester fosite.Requester) {
	// Merge implementation goes here
}

func (is *OidcSession) Sanitize(allowedParameters []string) fosite.Requester {
	// Sanitize implementation goes here
	return nil
}

func IDSessionOf(signature string, requester fosite.Requester) *OidcSession {
	jsonData, err := json.Marshal(requester.GetSession())
	if err != nil {
		log.Printf("Error occurred in FromRequester: %+v", err)
		return nil
	}

	return &OidcSession{
		Signature:         signature,
		RequestID:         requester.GetID(),
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

func (is *OidcSession) ToRequester() fosite.Requester {
	return &OidcSession{
		Signature:         is.Signature,
		RequestID:         is.RequestID,
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

func (is *OidcSession) ToAuthorizeRequest() fosite.AuthorizeRequester {
	ar := fosite.NewAuthorizeRequest()
	ar.Form = is.GetRequestForm()
	ar.RequestedAt = is.GetRequestedAt()
	ar.RequestedScope = is.GetRequestedScopes()
	ar.GrantedAudience = is.GetGrantedAudience()
	ar.GrantedScope = is.GetGrantedScopes()
	ar.Session = is.GetSession()
	ar.ID = is.GetID()

	return ar
}
