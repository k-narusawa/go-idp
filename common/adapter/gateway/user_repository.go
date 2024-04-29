package gateway

import (
	"log"

	"github.com/k-narusawa/go-idp/common/domain/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	db := Connect()
	var user models.User
	res := db.Where("username=?", username).First(&user)
	if res.Error != nil {
		log.Printf("Error occurred in FindByUsername: %+v", res.Error)
		return nil, res.Error
	}

	return &user, nil
}

func (r *UserRepository) Save(user *models.User) error {
	res := r.db.Save(user)
	if res.Error != nil {
		log.Printf("Error occurred in Save: %+v", res.Error)
		return res.Error
	}

	return nil
}
