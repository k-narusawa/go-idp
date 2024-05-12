package repository

import "github.com/k-narusawa/go-idp/domain/models"

type IAccessTokenRepository interface {
	FindBySubject(subject string) (*[]models.AccessToken, error)
	FindBySignature(signature string) (*models.AccessToken, error)
	DeleteBySignature(signature string) error
}
