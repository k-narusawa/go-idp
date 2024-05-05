package repository

import "github.com/k-narusawa/go-idp/authorization/domain/models"

type IClientRepository interface {
	FindClientByID(id string) (*models.Client, error)
	Save(client *models.Client) error
	DeleteByID(id string) error
}
