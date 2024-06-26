package gateways

import (
	"github.com/k-narusawa/go-idp/domain/models"
	"gorm.io/gorm"
)

type AccessTokenRepository struct {
	db *gorm.DB
}

func NewAccessTokenRepository(db *gorm.DB) *AccessTokenRepository {
	return &AccessTokenRepository{db}
}

func (r *AccessTokenRepository) FindBySubject(subject string) (*[]models.AccessToken, error) {
	var accessToken models.AccessToken
	res := r.db.Where("subject=?", subject).Find(&accessToken)
	if res.Error != nil {
		return nil, res.Error
	}

	return &[]models.AccessToken{accessToken}, nil
}

func (r *AccessTokenRepository) FindBySignature(signature string) (*models.AccessToken, error) {
	var accessToken models.AccessToken
	res := r.db.Where("signature=?", signature).Find(&accessToken)
	if res.Error != nil {
		return nil, res.Error
	}

	return &accessToken, nil
}

func (r *AccessTokenRepository) DeleteBySignature(signature string) error {
	res := r.db.Where("signature=?", signature).Delete(&models.AccessToken{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}
