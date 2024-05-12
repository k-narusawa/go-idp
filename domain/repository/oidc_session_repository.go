package repository

import "github.com/k-narusawa/go-idp/domain/models"

type IOidcSessionRepository interface {
	FindBySignature(signature string) (*models.OidcSession, error)
	DeleteBySignature(signature string) error
}
