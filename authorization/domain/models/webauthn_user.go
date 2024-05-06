package models

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type WebauthnUser struct {
	gorm.Model
	ID          string
	Name        string
	DisplayName string
	Credentials []WebauthnCredential
}

func NewWebauthnUser(userId string, userName string) *WebauthnUser {
	wu := &WebauthnUser{}
	wu.ID = userId
	wu.Name = userId
	wu.DisplayName = userName
	wu.Credentials = []WebauthnCredential{}
	return wu
}

func (wu WebauthnUser) WebAuthnID() []byte {
	return []byte(wu.ID)
}

func (wu WebauthnUser) WebAuthnName() string {
	return wu.Name
}

func (wu WebauthnUser) WebAuthnDisplayName() string {
	return wu.DisplayName
}

func (wu WebauthnUser) WebAuthnIcon() string {
	return ""
}

func (wu *WebauthnUser) AddCredential(cred webauthn.Credential) {
	wu.Credentials = append(wu.Credentials, *FromWebauthnCredential(wu.ID, &cred))
}

func (wu WebauthnUser) WebAuthnCredentials() []webauthn.Credential {
	credentials := []webauthn.Credential{}
	for _, cred := range wu.Credentials {
		credentials = append(credentials, *cred.To())
	}
	return credentials
}

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
