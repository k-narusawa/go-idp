package repository

import (
	"github.com/k-narusawa/go-idp/authorization/domain/models"
)

type ILoginSkipSessionRepository interface {
	FindByToken(token string) (*models.LoginSkipSession, error)
	Save(session *models.LoginSkipSession) error
	DeleteByToken(token string) error
}
