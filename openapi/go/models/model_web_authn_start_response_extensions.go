package models

type WebAuthnStartResponseExtensions struct {

	Appid string `json:"appid,omitempty"`

	AuthnSel string `json:"authnSel,omitempty"`

	Exts bool `json:"exts,omitempty"`
}
