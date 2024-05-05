package usecase

import (
	"log"
	"net/http"

	"github.com/k-narusawa/go-idp/authorization/domain/repository"
	"github.com/labstack/echo/v4"
)

type SessionUsecase struct {
	lssr repository.ILoginSkipSessionRepository
}

func NewSessionUsecase(
	lssr repository.ILoginSkipSessionRepository,
) SessionUsecase {
	return SessionUsecase{
		lssr: lssr,
	}
}

func (s *SessionUsecase) SetLoginSession(c echo.Context) error {
	log.Printf("SetSession")
	return c.Redirect(http.StatusFound, "/")
}
