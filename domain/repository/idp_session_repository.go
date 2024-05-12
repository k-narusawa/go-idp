package repository

import (
	"github.com/k-narusawa/go-idp/domain/models"
	"github.com/labstack/echo/v4"
)

type IIdpSessionRepository interface {
	Save(c echo.Context, idpSession *models.IdpSession) error
	Get(c echo.Context) (*models.IdpSession, error)
	Delete(c echo.Context) error
}
