package repository

import (
	"github.com/k-narusawa/go-idp/authorization/domain/models"
)

type IWebauthnCredentialRepository interface {
	FindByID(id string) (*models.WebauthnCredential, error)
	FindByUserID(userID string) ([]models.WebauthnCredential, error)
	Save(credential *models.WebauthnCredential) error
	DeleteByCredentialID(credentialID uint) error
}
