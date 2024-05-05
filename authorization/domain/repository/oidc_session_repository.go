package repository

import "github.com/k-narusawa/go-idp/authorization/domain/models"

type IOidcSessionRepository interface {
	FindBySignature(signature string) (*models.OidcSession, error)
	DeleteBySignature(signature string) error
}
