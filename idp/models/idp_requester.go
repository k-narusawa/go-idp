package models

import (
	"net/url"
	"time"

	"github.com/ory/fosite"
	"gorm.io/gorm"
)

type IdpRequester struct {
	gorm.Model
	ID                string           `gorm:"type:varchar(255);not null;unique" json:"id"`
	RequestedAt       time.Time        `json:"requested_at"`
	Client            Client           `json:"client"`
	RequestedScopes   fosite.Arguments `gorm:"type:text" json:"requested_scopes"`   // JSON-encoded []string
	RequestedAudience fosite.Arguments `gorm:"type:text" json:"requested_audience"` // JSON-encoded []string
	GrantedScopes     fosite.Arguments `gorm:"type:text" json:"granted_scopes"`     // JSON-encoded []string
	GrantedAudience   fosite.Arguments `gorm:"type:text" json:"granted_audience"`   // JSON-encoded []string
	Session           string           `json:"session"`
	RequestForm       url.Values       `json:"request_form"`
}

func (r *IdpRequester) SetID(id string) {
	r.ID = id
}

func (r *IdpRequester) GetID() string {
	return r.ID
}

func (r *IdpRequester) GetRequestedAt() time.Time {
	return r.RequestedAt
}

func (r *IdpRequester) GetClient() fosite.Client {
	return CastToClient(r.Client)
}

func (r *IdpRequester) GetRequestedScopes() fosite.Arguments {
	return r.RequestedScopes
}

func (r *IdpRequester) GetRequestedAudience() fosite.Arguments {
	return r.RequestedAudience
}

func (r *IdpRequester) SetRequestedScopes(scopes fosite.Arguments) {
	r.RequestedScopes = scopes
}

func (r *IdpRequester) SetRequestedAudience(audience fosite.Arguments) {
	r.RequestedAudience = audience
}

func (r *IdpRequester) AppendRequestedScope(scope string) {
	r.RequestedScopes = append(r.RequestedScopes, scope)
}

func (r *IdpRequester) GetGrantedScopes() fosite.Arguments {
	return r.GrantedScopes
}

func (r *IdpRequester) GetGrantedAudience() fosite.Arguments {
	return r.GrantedAudience
}

func (r *IdpRequester) GrantScope(scope string) {
	r.GrantedScopes = append(r.GrantedScopes, scope)
}

func (r *IdpRequester) GrantAudience(audience string) {
	r.GrantedAudience = append(r.GrantedAudience, audience)
}

func (r *IdpRequester) GetSession() fosite.Session {
	return nil
}

func (r *IdpRequester) SetSession(session fosite.Session) {
	// r.Session = nil
}

func (r *IdpRequester) GetRequestForm() url.Values {
	return r.RequestForm
}

func (r *IdpRequester) Merge(requester fosite.Requester) {
	// Implement this method based on your requirements
}

func (r *IdpRequester) Sanitize(allowedParameters []string) IdpRequester {
	// Implement this method based on your requirements
	return IdpRequester{}
}
