package gateways

import (
	"github.com/k-narusawa/go-idp/domain/models"
	"gorm.io/gorm"
)

type WebauthnSessionRepository struct {
	db *gorm.DB
}

func NewWebauthnSessionRepository(db *gorm.DB) *WebauthnSessionRepository {
	return &WebauthnSessionRepository{db: db}
}

func (r *WebauthnSessionRepository) FindByChallenge(challenge string) (*models.WebauthnSessionData, error) {
	wsd := models.WebauthnSessionData{}
	result := r.db.Where("challenge = ?", challenge).First(&wsd)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wsd, nil
}

func (r *WebauthnSessionRepository) Save(credential *models.WebauthnSessionData) error {
	result := r.db.Create(credential)
	return result.Error
}

func (r *WebauthnSessionRepository) DeleteByChallenge(challenge string) error {
	result := r.db.Where("challenge = ?", challenge).Delete(&models.WebauthnSessionData{})
	return result.Error
}
