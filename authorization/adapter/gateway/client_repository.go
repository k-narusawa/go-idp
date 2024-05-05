package gateway

import (
	"github.com/k-narusawa/go-idp/authorization/domain/models"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{db}
}

func (r *ClientRepository) FindClientByID(id string) (*models.Client, error) {
	var client models.Client
	res := r.db.Where("id=?", id).First(&client)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, res.Error
	}

	return &client, nil
}

func (r *ClientRepository) Save(client *models.Client) error {
	res := r.db.Save(client)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *ClientRepository) DeleteByID(id string) error {
	res := r.db.Where("id=?", id).Delete(&models.Client{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}
