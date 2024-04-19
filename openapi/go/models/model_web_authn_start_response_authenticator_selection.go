package models

type WebAuthnStartResponseAuthenticatorSelection struct {

	AuthenticatorAttachment string `json:"authenticatorAttachment,omitempty"`

	RequireResidentKey bool `json:"requireResidentKey,omitempty"`

	UserVerification string `json:"userVerification,omitempty"`
}
