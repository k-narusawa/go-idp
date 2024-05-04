package models

type ClientRequest struct {

	// クライアントID
	ClientId string `json:"client_id"`

	// クライアントシークレット
	ClientSecret string `json:"client_secret"`

	RedirectUris []string `json:"redirect_uris"`

	GrantTypes []string `json:"grant_types"`

	ResponseTypes []string `json:"response_types"`

	// クライアントがサポートするスコープ - openid: OpenID Connectのスコープ - offline: リフレッシュトークンを取得するためのスコープ 
	Scopes []string `json:"scopes"`
}
