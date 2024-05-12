package application

import (
	"strings"

	"github.com/k-narusawa/go-idp/domain/models"
	"github.com/k-narusawa/go-idp/domain/repository"
	"github.com/labstack/echo/v4"
)

type ClientInteractor struct {
	cr repository.IClientRepository
}

func NewClientInteractor(
	cr repository.IClientRepository,
) ClientInteractor {
	return ClientInteractor{
		cr: cr,
	}
}

func (cu ClientInteractor) Register(c echo.Context) error {
	req := new(ClientRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	client := models.NewClient(
		req.ClientId,
		req.ClientSecret, req.RedirectUris, req.GrantTypes, req.ResponseTypes, req.Scopes,
		req.Audience,
		req.Public,
	)

	if err := cu.cr.Save(client); err != nil {
		return err
	}

	res := &ClientResponse{
		ClientId:      client.ID,
		RedirectUris:  strings.Split(client.RedirectURIs, ","),
		ResponseTypes: strings.Split(client.ResponseTypes, ","),
		GrantTypes:    strings.Split(client.GrantTypes, ","),
		Scopes:        strings.Split(client.Scopes, ","),
		Audience:      client.Audience,
		Public:        client.Public,
	}

	return c.JSON(201, res)
}

func (cu ClientInteractor) Get(c echo.Context) error {
	clientId := c.Param("id")

	client, err := cu.cr.FindClientByID(clientId)
	if err != nil {
		return err
	}

	if client == nil {
		return c.NoContent(404)
	}

	res := &ClientResponse{
		ClientId:      client.ID,
		RedirectUris:  strings.Split(client.RedirectURIs, ","),
		ResponseTypes: strings.Split(client.ResponseTypes, ","),
		GrantTypes:    strings.Split(client.GrantTypes, ","),
		Scopes:        strings.Split(client.Scopes, ","),
		Audience:      client.Audience,
		Public:        client.Public,
	}

	return c.JSON(200, res)
}

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

	// オーディエンス
	Audience string `json:"audience,omitempty"`

	// ClientSecretを安全に管理できるかどうか
	Public bool `json:"public,omitempty"`
}

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
