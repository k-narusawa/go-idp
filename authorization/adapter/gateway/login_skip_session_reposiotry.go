package gateway

import (
	"github.com/k-narusawa/go-idp/authorization/domain/models"
	"gorm.io/gorm"
)

type LoginSkipSessionRepository struct {
	db *gorm.DB
}

func NewLoginSkipSessionRepository(db *gorm.DB) *LoginSkipSessionRepository {
	return &LoginSkipSessionRepository{
		db: db,
	}
}

func (r *LoginSkipSessionRepository) Save(session *models.LoginSkipSession) error {
	result := r.db.Create(session)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *LoginSkipSessionRepository) FindByToken(token string) (*models.LoginSkipSession, error) {
	var lss models.LoginSkipSession

	result := r.db.Where("token=?", token).First(&lss)

	if result.Error != nil {
		return nil, result.Error
	}

	return &lss, nil
}

func (r *LoginSkipSessionRepository) DeleteByToken(token string) error {
	res := r.db.Where("token=?", token).Delete(&models.LoginSkipSession{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}
