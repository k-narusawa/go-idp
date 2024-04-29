package repository

import (
	"github.com/k-narusawa/go-idp/authorization/domain/models"

	"github.com/labstack/echo/v4"
)

type IIdSessionRepository interface {
	CreateIdSession(c echo.Context, is models.IDSession) error
	GetIdSession(c echo.Context) (*models.IDSession, error)
}
