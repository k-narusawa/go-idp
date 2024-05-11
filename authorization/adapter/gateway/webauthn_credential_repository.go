package gateway

import (
	"github.com/k-narusawa/go-idp/authorization/domain/models"
	"gorm.io/gorm"
)

type WebauthnCredentialRepository struct {
	db *gorm.DB
}

func NewWebauthnCredentialRepository(db *gorm.DB) *WebauthnCredentialRepository {
	return &WebauthnCredentialRepository{
		db: db,
	}
}

func (r *WebauthnCredentialRepository) FindByID(id string) (*models.WebauthnCredential, error) {
	credential := models.WebauthnCredential{}
	result := r.db.Where("id = ?", id).First(&credential)
	if result.Error != nil {
		return nil, result.Error
	}
	return &credential, nil
}

func (r *WebauthnCredentialRepository) FindByUserID(userID string) ([]models.WebauthnCredential, error) {
	credentials := []models.WebauthnCredential{}
	result := r.db.Where("user_id = ?", userID).Find(&credentials)
	if result.Error != nil {
		return nil, result.Error
	}
	return credentials, nil
}

func (r *WebauthnCredentialRepository) Save(credential *models.WebauthnCredential) error {
	result := r.db.Create(credential)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *WebauthnCredentialRepository) DeleteByCredentialID(credentialID uint) error {
	result := r.db.Where("credential_id = ?", credentialID).Delete(&models.WebauthnCredential{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
