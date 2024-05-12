package application

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/k-narusawa/go-idp/domain/models"
	"github.com/k-narusawa/go-idp/domain/repository"
	"github.com/k-narusawa/go-idp/logger"
)

type IntrospectInteractor struct {
	logger logger.Logger
	atr    repository.IAccessTokenRepository
}

func NewIntrospectInteractor(
	logger logger.Logger,
	atr repository.IAccessTokenRepository,
) IntrospectInteractor {
	return IntrospectInteractor{
		logger: logger,
		atr:    atr,
	}
}

func (i *IntrospectInteractor) Introspect(token string) (accessToken *models.AccessToken, err error) {
	i.logger.Info("Introspect", "token", token)
	splited := strings.Split(token, ".")
	if len(splited) != 2 {
		slog.Info("Failed to introspect token: invalid token")
		return
	}

	signature := splited[1]

	accessToken, err = i.atr.FindBySignature(signature)
	if err != nil {
		i.logger.Warn("Failed to introspect", "err", err)
		return nil, err
	}

	if accessToken.IsExpired() {
		i.logger.Warn("Failed to introspect", "err", "token is expired")
		i.atr.DeleteBySignature(signature)
		return nil, errors.New("token is expired")
	}

	// TODO: scopeのチェックとか
	i.logger.Info("GrantedScope", "scopes", accessToken.GrantedScope)

	return accessToken, nil
}
