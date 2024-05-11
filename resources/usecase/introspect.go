package usecase

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/k-narusawa/go-idp/authorization/domain/models"
	"github.com/k-narusawa/go-idp/authorization/domain/repository"
)

type IntrospectUsecase struct {
	atr repository.IAccessTokenRepository
}

func NewIntrospectUsecase(atr repository.IAccessTokenRepository) IntrospectUsecase {
	return IntrospectUsecase{
		atr: atr,
	}
}

func (i *IntrospectUsecase) Introspect(token string) (accessToken *models.AccessToken, err error) {
	slog.Info("Introspect", "tokent", token)
	splited := strings.Split(token, ".")
	if len(splited) != 2 {
		slog.Info("Failed to introspect token: invalid token")
		return
	}

	signature := splited[1]

	accessToken, err = i.atr.FindBySignature(signature)
	if err != nil {
		slog.Warn("Failed to introspect", "err", err)
		return nil, err
	}

	if accessToken.IsExpired() {
		slog.Warn("Failed to introspect", "err", "token is expired")
		i.atr.DeleteBySignature(signature)
		return nil, errors.New("token is expired")
	}

	// TODO: scopeのチェックとか
	slog.Info("GrantedScope", "scopes", accessToken.GrantedScope)

	return accessToken, nil
}
