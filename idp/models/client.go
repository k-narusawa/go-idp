package models

import (
	"log"
	"strings"

	"github.com/ory/fosite"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ID             string `gorm:"type:varchar(255);not null;unique" json:"id"`
	Secret         []byte `gorm:"type:blob" json:"client_secret,omitempty"`
	RotatedSecrets string `gorm:"type:text" json:"rotated_secrets,omitempty"` // JSON-encoded [][]byte
	RedirectURIs   string `gorm:"type:text" json:"redirect_uris"`             // JSON-encoded []string
	GrantTypes     string `gorm:"type:text" json:"grant_types"`               // JSON-encoded []string
	ResponseTypes  string `gorm:"type:text" json:"response_types"`            // JSON-encoded []string
	Scopes         string `gorm:"type:text" json:"scopes"`                    // JSON-encoded []string
	Audience       string `gorm:"type:text" json:"audience"`                  // JSON-encoded []string
	Public         bool   `gorm:"type:boolean" json:"public"`
}

func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) GetHashedSecret() []byte {
	return c.Secret
}

func (c *Client) GetRedirectURIs() []string {
	return []string{"http://localhost:3846/callback"}
}

func (c *Client) GetGrantTypes() fosite.Arguments {
	log.Printf("grant_types: %+v", strings.Split(c.GrantTypes, ","))
	// arr := strings.Split(c.GrantTypes, ",")
	return fosite.Arguments{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"}
}

func (c *Client) GetResponseTypes() fosite.Arguments {
	log.Printf("response_types: %+v", c.ResponseTypes)
	// return fosite.Arguments(strings.Split(c.ResponseTypes, ","))
	// return fosite.Arguments{"id_token", "code", "token", "id_token token", "code id_token", "code token", "code id_token token"}
	return fosite.Arguments{"code"}
}

func (c *Client) GetScopes() fosite.Arguments {
	log.Printf("scopes: %+v", strings.Split(c.Scopes, ","))
	return fosite.Arguments((strings.Split(c.Scopes, ",")))
}

func (c *Client) IsPublic() bool {
	log.Printf("public: %+v", c.Public)
	return c.Public
}

func (c *Client) GetAudience() fosite.Arguments {
	log.Printf("audience: %+v", strings.Split(c.Audience, ","))
	return fosite.Arguments{}
}

func CastToClient(mc Client) fosite.Client {
	return &Client{
		ID:             mc.ID,
		Secret:         mc.Secret,
		RotatedSecrets: mc.RotatedSecrets,
		RedirectURIs:   mc.RedirectURIs,
		ResponseTypes:  mc.ResponseTypes,
		GrantTypes:     mc.GrantTypes,
		Scopes:         mc.Scopes,
		Audience:       mc.Audience,
		Public:         mc.Public,
	}
}
