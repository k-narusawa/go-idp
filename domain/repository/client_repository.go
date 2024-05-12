package repository

import "github.com/k-narusawa/go-idp/domain/models"

type IClientRepository interface {
	FindClientByID(id string) (*models.Client, error)
	Save(client *models.Client) error
	DeleteByID(id string) error
}
