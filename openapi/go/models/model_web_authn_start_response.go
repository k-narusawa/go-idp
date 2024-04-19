package models

type WebAuthnStartResponse struct {

	Rp WebAuthnStartResponseRp `json:"rp"`

	User WebAuthnStartResponseUser `json:"user"`

	Challenge string `json:"challenge"`

	PubKeyCredParams []WebAuthnStartResponsePubKeyCredParamsInner `json:"pubKeyCredParams"`

	Timeout int32 `json:"timeout"`

	ExcludeCredentials []WebAuthnStartResponseExcludeCredentialsInner `json:"excludeCredentials"`

	AuthenticatorSelection WebAuthnStartResponseAuthenticatorSelection `json:"authenticatorSelection"`

	Attestation string `json:"attestation"`

	Extensions WebAuthnStartResponseExtensions `json:"extensions,omitempty"`
}
