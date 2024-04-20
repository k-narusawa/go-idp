package model

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebauthnUser struct {
	id          uint64                `gorm:"primary_key"`
	name        string                `gorm:"unique"`
	displayName string                `gorm:"unique"`
	credentials []webauthn.Credential `gorm:"many2many:webauthn_user_credentials"`
}

func NewUser(name string, displayName string) *WebauthnUser {
	user := &WebauthnUser{}
	user.id = randomUint64()
	user.name = name
	user.displayName = displayName
	return user
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// WebAuthnID returns the user's ID
func (wu WebauthnUser) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(wu.id))
	return buf
}

// WebAuthnName returns the user's username
func (wu WebauthnUser) WebAuthnName() string {
	return wu.name
}

// WebAuthnDisplayName returns the user's display name
func (wu WebauthnUser) WebAuthnDisplayName() string {
	return wu.displayName
}

// WebAuthnIcon is not (yet) implemented
func (wu WebauthnUser) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (wu *WebauthnUser) AddCredential(cred webauthn.Credential) {
	wu.credentials = append(wu.credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (wu WebauthnUser) WebAuthnCredentials() []webauthn.Credential {
	return wu.credentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (wu WebauthnUser) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range wu.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}
