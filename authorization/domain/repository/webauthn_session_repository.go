package repository

import (
	"github.com/k-narusawa/go-idp/authorization/domain/models"
)

type IWebauthnSessionRepository interface {
	FindByChallenge(challenge string) (*models.WebauthnSessionData, error)
	Save(credential *models.WebauthnSessionData) error
	DeleteByChallenge(challenge string) error
}
