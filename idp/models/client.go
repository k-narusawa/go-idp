package models

import (
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
	return strings.Split(c.GrantTypes, ",")
}

func (c *Client) GetResponseTypes() fosite.Arguments {
	return fosite.Arguments(strings.Split(c.ResponseTypes, ","))
}

func (c *Client) GetScopes() fosite.Arguments {
	return fosite.Arguments((strings.Split(c.Scopes, ",")))
}

func (c *Client) IsPublic() bool {
	return c.Public
}

func (c *Client) GetAudience() fosite.Arguments {
	return strings.Split(c.Audience, ",")
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

func ClientOf(fc fosite.Client) Client {
	c := fc.(*Client)
	return Client{
		ID:             c.ID,
		Secret:         c.Secret,
		RotatedSecrets: c.RotatedSecrets,
		RedirectURIs:   c.RedirectURIs,
		ResponseTypes:  c.ResponseTypes,
		GrantTypes:     c.GrantTypes,
		Scopes:         c.Scopes,
		Audience:       c.Audience,
		Public:         c.Public,
	}
}
