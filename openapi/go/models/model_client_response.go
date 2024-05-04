package models

type ClientResponse struct {

	// クライアントID
	ClientId string `json:"client_id"`

	RedirectUris []string `json:"redirect_uris"`

	GrantTypes []string `json:"grant_types"`

	ResponseTypes []string `json:"response_types"`

	// クライアントがサポートするスコープ - openid: OpenID Connectのスコープ - offline: リフレッシュトークンを取得するためのスコープ 
	Scopes []string `json:"scopes"`

	// オーディエンス 
	Audience string `json:"audience,omitempty"`

	// ClientSecretを安全に管理できるかどうか 
	Public bool `json:"public,omitempty"`
}
