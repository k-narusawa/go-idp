package usecase

import (
	"errors"
	"log"
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
	splited := strings.Split(token, ".")
	if len(splited) != 2 {
		log.Printf("Failed to introspect token: invalid token")
		return
	}

	signature := splited[1]

	log.Printf("Introspecting token: %s", signature)

	accessToken, err = i.atr.FindBySignature(signature)
	if err != nil {
		log.Printf("Failed to introspect token: %v", err)
		return nil, err
	}

	if accessToken.IsExpired() {
		log.Printf("Failed to introspect token: token is expired")
		i.atr.DeleteBySignature(signature)
		return nil, errors.New("token is expired")
	}

	// TODO: scopeのチェックとか
	log.Printf("GrantedScope %v", accessToken.GetGrantedScopes())

	return accessToken, nil
}
