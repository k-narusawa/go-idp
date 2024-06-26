package repository

import (
	"github.com/k-narusawa/go-idp/domain/models"
)

type ILoginSkipSessionRepository interface {
	FindByToken(token string) (*models.LoginSkipSession, error)
	Save(session *models.LoginSkipSession) error
	DeleteByToken(token string) error
}
