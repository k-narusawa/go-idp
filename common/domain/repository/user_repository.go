package repository

import "github.com/k-narusawa/go-idp/common/domain/models"

type IUserRepository interface {
	FindByUsername(username string) (*models.User, error)
	Save(user *models.User) error
}
