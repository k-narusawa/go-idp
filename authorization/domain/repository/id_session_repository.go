package repository

import (
	"idp/authorization/domain/models"

	"github.com/labstack/echo/v4"
)

type IIdSessionRepository interface {
	CreateIdSession(c echo.Context, is models.IDSession) error
	GetIdSession(c echo.Context) (*models.IDSession, error)
}
