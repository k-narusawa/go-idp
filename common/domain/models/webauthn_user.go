package models

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type WebauthnUser struct {
	gorm.Model
	WuID        int64                `gorm:"primaryKey"`
	ID          string               `gorm:"type:text"`
	Name        string               `gorm:"varchar(255)"`
	DisplayName string               `gorm:"varchar(255)"`
	Credentials []WebauthnCredential `gorm:"foreignKey:ID;references:ID"`
}

func NewWebauthnUser(name string, displayName string) *WebauthnUser {
	user := &WebauthnUser{}
	user.ID = name
	user.Name = name
	user.DisplayName = displayName
	return user
}

// WebAuthnID returns the user's ID
func (wu WebauthnUser) WebAuthnID() []byte {
	return []byte(wu.ID)
}

// WebAuthnName returns the user's username
func (wu WebauthnUser) WebAuthnName() string {
	return wu.Name
}

// WebAuthnDisplayName returns the user's display name
func (wu WebauthnUser) WebAuthnDisplayName() string {
	return wu.DisplayName
}

// WebAuthnIcon is not (yet) implemented
func (wu WebauthnUser) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (wu *WebauthnUser) AddCredential(cred webauthn.Credential) {
	wu.Credentials = append(wu.Credentials, *FromWebauthnCredential(&cred))
}

// WebAuthnCredentials returns credentials owned by the user
func (wu WebauthnUser) WebAuthnCredentials() []webauthn.Credential {
	credentials := []webauthn.Credential{}
	for _, cred := range wu.Credentials {
		credentials = append(credentials, *cred.ToWebauthnCredential())
	}
	return credentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (wu WebauthnUser) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	// for _, cred := range wu.Credentials {
	// 	descriptor := protocol.CredentialDescriptor{
	// 		Type:         protocol.PublicKeyCredentialType,
	// 		CredentialID: cred.ID,
	// 	}
	// 	credentialExcludeList = append(credentialExcludeList, descriptor)
	// }

	return credentialExcludeList
}
