package gateway

import (
	"errors"
	"log"

	"github.com/k-narusawa/go-idp/domain/models"
	"github.com/ory/fosite"
	"gorm.io/gorm"
)

type OidcSessionRepository struct {
	db *gorm.DB
}

func NewOidcSessionRepository(db *gorm.DB) *OidcSessionRepository {
	return &OidcSessionRepository{
		db: db,
	}
}

func (r *OidcSessionRepository) FindBySignature(signature string) (*models.OidcSession, error) {
	var is models.OidcSession

	result := r.db.Preload("Client").Where("signature=?", signature).First(&is)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No record found for signature: %s", signature)
			return nil, fosite.ErrNotFound
		}
		log.Printf("Error occurred in GetOpenIDConnectSession: %+v", result.Error)
		return nil, result.Error
	}

	return &is, nil
}

func (r *OidcSessionRepository) DeleteBySignature(signature string) error {
	res := r.db.Where("signature=?", signature).Delete(&models.OidcSession{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}
