package usecase

import (
	"strings"

	"github.com/k-narusawa/go-idp/authorization/domain/models"
	"github.com/k-narusawa/go-idp/authorization/domain/repository"
	"github.com/labstack/echo/v4"
)

type ClientUsecase struct {
	cr repository.IClientRepository
}

func NewClientUsecase(
	cr repository.IClientRepository,
) ClientUsecase {
	return ClientUsecase{
		cr: cr,
	}
}

func (cu ClientUsecase) Register(c echo.Context) error {
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
