package model

type IntrospectResponse struct {
	Active   bool     `json:"active"`
	Aud      []string `json:"aud"`
	ClientId string   `json:"client_id"`
	Exp      int      `json:"exp"`
	Iat      int      `json:"iat"`
	Scope    string   `json:"scope"`
	Sub      string   `json:"sub"`
}
