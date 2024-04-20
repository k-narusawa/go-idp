package model

import (
	"strings"
	"time"

	"encoding/json"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type WebauthnSessionData struct {
	gorm.Model
	Challenge            string    `gorm:"unique" json:"challenge"`
	UserID               []byte    `json:"user_id"`
	AllowedCredentialIDs string    `json:"allowed_credentials,omitempty"`
	Expires              time.Time `json:"expires"`
	UserVerification     string    `json:"userVerification"`
	Extensions           string    `json:"extensions,omitempty"`
}

func FromSessionData(sd *webauthn.SessionData) *WebauthnSessionData {
	acIds := make([]string, len(sd.AllowedCredentialIDs))
	ej, err := json.Marshal(sd.Extensions)
	if err != nil {
		panic(err)
	}

	return &WebauthnSessionData{
		Challenge:            sd.Challenge,
		UserID:               sd.UserID,
		AllowedCredentialIDs: strings.Join(acIds, ","),
		Expires:              sd.Expires,
		UserVerification:     string(sd.UserVerification),
		Extensions:           string(ej),
	}
}

func (wsd *WebauthnSessionData) ToSessionData() *webauthn.SessionData {
	allowedCredentialIDs := make([][]byte, len(strings.Split(wsd.AllowedCredentialIDs, ",")))
	uv := protocol.UserVerificationRequirement(wsd.UserVerification)
	var ex map[string]interface{}
	if err := json.Unmarshal([]byte(wsd.Extensions), &ex); err != nil {
		panic(err)
	}

	return &webauthn.SessionData{
		Challenge:            wsd.Challenge,
		UserID:               wsd.UserID,
		AllowedCredentialIDs: allowedCredentialIDs,
		Expires:              wsd.Expires,
		UserVerification:     uv,
		Extensions:           protocol.AuthenticationExtensions(ex),
	}
}
