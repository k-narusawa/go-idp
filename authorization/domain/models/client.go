package models

import (
	"strings"

	"github.com/ory/fosite"
)

type Client struct {
	ID            string `gorm:"type:varchar(255);not null;unique" `
	Secret        []byte `gorm:"type:blob"`
	RedirectURIs  string `gorm:"type:text"`
	GrantTypes    string `gorm:"type:text"`
	ResponseTypes string `gorm:"type:text"`
	Scopes        string `gorm:"type:text"`
	Audience      string `gorm:"type:text"`
	Public        bool   `gorm:"type:boolean"`
}

func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) GetHashedSecret() []byte {
	return c.Secret
}

func (c *Client) GetRedirectURIs() []string {
	return strings.Split(c.RedirectURIs, ",")
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
		ID:            mc.ID,
		Secret:        mc.Secret,
		RedirectURIs:  mc.RedirectURIs,
		ResponseTypes: mc.ResponseTypes,
		GrantTypes:    mc.GrantTypes,
		Scopes:        mc.Scopes,
		Audience:      mc.Audience,
		Public:        mc.Public,
	}
}

func ClientOf(fc fosite.Client) Client {
	c := fc.(*Client)
	return Client{
		ID:            c.ID,
		Secret:        c.Secret,
		RedirectURIs:  c.RedirectURIs,
		ResponseTypes: c.ResponseTypes,
		GrantTypes:    c.GrantTypes,
		Scopes:        c.Scopes,
		Audience:      c.Audience,
		Public:        c.Public,
	}
}
