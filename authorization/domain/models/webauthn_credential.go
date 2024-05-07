package models

import (
	"encoding/json"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebauthnCredential struct {
	ID              []byte `gorm:"type:blob;index;primary_key"`
	UserID          string `gorm:"type:varchar(36);not null;index"`
	PublicKey       []byte `gorm:"type:blob"`
	AttestationType string `gorm:"type:text"`
	Transport       string `gorm:"type:text"`
	Flags           string `gorm:"type:text"`
	Authenticator   string `gorm:"type:text"`
}

func (c *WebauthnCredential) To() *webauthn.Credential {
	tp := strings.Split(c.Transport, ",")
	var transport []protocol.AuthenticatorTransport
	for _, t := range tp {
		transport = append(transport, protocol.AuthenticatorTransport(t))
	}

	var cf webauthn.CredentialFlags
	err := json.Unmarshal([]byte(c.Flags), &cf)
	if err != nil {
		panic(err)
	}

	a := webauthn.Authenticator{}
	if c.Authenticator != "" {
		err = json.Unmarshal([]byte(c.Authenticator), &a)
		if err != nil {
			panic(err)
		}
	}

	return &webauthn.Credential{
		ID:              c.ID,
		PublicKey:       c.PublicKey,
		AttestationType: c.AttestationType,
		Transport:       transport,
		Flags:           cf,
		Authenticator:   a,
	}
}

func FromWebauthnCredential(userId string, cred *webauthn.Credential) *WebauthnCredential {
	var transport []string
	for _, t := range cred.Transport {
		transport = append(transport, string(t))
	}

	cf, err := json.Marshal(cred.Flags)
	if err != nil {
		panic(err)
	}

	a, err := json.Marshal(cred.Authenticator)
	if err != nil {
		panic(err)
	}

	return &WebauthnCredential{
		ID:              cred.ID,
		UserID:          userId,
		PublicKey:       cred.PublicKey,
		AttestationType: cred.AttestationType,
		Transport:       strings.Join(transport, ","),
		Flags:           string(cf),
		Authenticator:   string(a),
	}
}
