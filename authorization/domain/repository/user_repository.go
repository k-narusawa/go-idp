package repository

import "github.com/k-narusawa/go-idp/authorization/domain/models"

type IUserRepository interface {
	FindByUsername(username string) (*models.User, error)
	FindByUserID(userID string) (*models.User, error)
	Save(user *models.User) error
}
