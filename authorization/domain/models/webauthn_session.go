package models

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
	// sd.AllowedCredentialIDsという[][]byteを[]stringに変換
	acIds := make([]string, len(sd.AllowedCredentialIDs))
	for i := range sd.AllowedCredentialIDs {
		acIds[i] = string(sd.AllowedCredentialIDs[i])
	}

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
	for i, id := range strings.Split(wsd.AllowedCredentialIDs, ",") {
		allowedCredentialIDs[i] = []byte(id)
	}
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
