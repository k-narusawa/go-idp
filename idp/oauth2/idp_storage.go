package oauth2

import (
	"context"
	"errors"
	"log"
	"time"

	"idp/infrastructure"
	"idp/models"

	"github.com/ory/fosite"
	"gorm.io/gorm"
)

type IdpStorage struct {
	Clients            []models.Client
	Users              []models.User
	IDSessions         []models.IDSession
	AuthorizationCodes []models.AuthorizationCode
	AccessTokens       []models.AccessToken
	RefreshTokens      []models.RefreshToken
}

func NewIdpStorage() *IdpStorage {
	is := IdpStorage{}
	return &is
}

func (s *IdpStorage) CreateClient(_ context.Context, client fosite.Client) {
	db := infrastructure.Connect()

	c := models.ClientOf(client)
	result := db.Create(&c)
	if result.Error != nil {
		log.Printf("Error occurred in CreateClient: %+v", result.Error)
	}
	log.Printf("CreateClient ClientID: %+v", c.ID)
}

func (s *IdpStorage) GetClient(_ context.Context, id string) (fosite.Client, error) {
	db := infrastructure.Connect()

	var c models.Client
	res := db.Where("id=?", id).First(&c)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Printf("No record found for id: %s", id)
			return nil, fosite.ErrNotFound
		}
		log.Printf("Error occurred in GetClient: %+v", res.Error)
		return nil, res.Error
	}

	log.Printf("GetClient ClientID: %+v", c.ID)

	return models.CastToClient(c), nil
}

func (s *IdpStorage) ClientAssertionJWTValid(_ context.Context, jti string) error {
	log.Printf("ClientAssertionJWTValid: %+v", jti)
	return nil
}

func (s *IdpStorage) SetClientAssertionJWT(_ context.Context, jti string, exp time.Time) error {
	log.Printf("ClientAssertionJWTValid: %+v", jti)
	return nil
}

func (s *IdpStorage) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	db := infrastructure.Connect()

	at := models.AccessTokenOf(signature, request)

	result := db.Create(&at)

	if result.Error != nil {
		log.Printf("Error occurred in CreateAccessTokenSession: %+v", result.Error)
		return result.Error
	}

	log.Printf("CreateAccessTokenSession Signature: %+v", signature)
	return nil
}

func (s *IdpStorage) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	db := infrastructure.Connect()

	var at models.AccessToken
	result := db.Where("signature=?", signature).First(&at)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No record found for signature: %s", signature)
			return nil, fosite.ErrNotFound
		}
		log.Printf("Error occurred in GetAccessTokenSession: %+v", result.Error)
		return nil, result.Error
	}

	log.Printf("GetAccessTokenSession Signature: %+v", signature)
	return at.ToRequester(), nil
}

func (s *IdpStorage) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	db := infrastructure.Connect()

	result := db.Where("signature=?", signature).Delete(&models.AccessToken{})
	if result.Error != nil {
		log.Printf("Error occurred in DeleteAccessTokenSession: %+v", result.Error)
		return result.Error
	}

	log.Printf("DeleteAccessTokenSession Signature: %+v", signature)
	return nil
}

func (s *IdpStorage) CreateAuthorizeCodeSession(_ context.Context, code string, req fosite.Requester) error {
	db := infrastructure.Connect()

	ac := models.AuthorizationCodeOf(code, req)

	result := db.Create(&ac)
	if result.Error != nil {
		log.Printf("Error occurred in CreateAuthorizeCodeSession: %+v", result.Error)
		return result.Error
	}

	log.Printf("CreateAuthorizeCodeSession Code: %+v", code)
	return nil
}

func (s *IdpStorage) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {
	db := infrastructure.Connect()

	var ac models.AuthorizationCode
	ar := db.
		Preload("Client").
		Where("signature=?", code).
		Find(&ac)

	if ar.Error != nil {
		log.Printf("Error occurred in GetAuthorizeCodeSession: %+v", ar.Error)
		return nil, ar.Error
	}

	log.Printf("GetAuthorizeCodeSession Code: %+v", code)
	return ac.ToRequester(), nil
}

func (s *IdpStorage) InvalidateAuthorizeCodeSession(ctx context.Context, code string) (err error) {
	db := infrastructure.Connect()

	result := db.Where("signature=?", code).Delete(&models.AuthorizationCode{})
	if result.Error != nil {
		log.Printf("Error occurred in InvalidateAuthorizeCodeSession: %+v", result.Error)
		return result.Error
	}

	log.Printf("InvalidateAuthorizeCodeSession Code: %+v", code)
	return nil
}

func (s *IdpStorage) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	db := infrastructure.Connect()

	rt := models.RefreshTokenOf(signature, request)
	result := db.Create(&rt)

	if result.Error != nil {
		log.Printf("Error occurred in CreateRefreshTokenSession: %+v", result.Error)
		return result.Error
	}

	log.Printf("CreateRefreshTokenSession Signature: %+v", signature)
	return nil
}

func (s *IdpStorage) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	db := infrastructure.Connect()

	var rt models.RefreshToken
	result := db.Where("signature=?", signature).First(&rt)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No record found for signature: %s", signature)
			return nil, fosite.ErrNotFound
		}
		log.Printf("Error occurred in GetRefreshTokenSession: %+v", result.Error)
		return nil, result.Error
	}

	log.Printf("GetRefreshTokenSession Signature: %+v", signature)
	return rt.ToRequester(), nil
}

func (s *IdpStorage) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	db := infrastructure.Connect()

	result := db.Where("signature=?", signature).Delete(&models.RefreshToken{})
	if result.Error != nil {
		log.Printf("Error occurred in DeleteRefreshTokenSession: %+v", result.Error)
		return result.Error
	}

	log.Printf("DeleteRefreshTokenSession Signature: %+v", signature)
	return nil
}

func (s *IdpStorage) RevokeAccessToken(ctx context.Context, requestID string) error {
	db := infrastructure.Connect()

	result := db.Where("signature=?", requestID).Delete(&models.AccessToken{})
	if result.Error != nil {
		log.Printf("Error occurred in RevokeAccessToken: %+v", result.Error)
		return result.Error
	}

	log.Printf("RevokeAccessToken RequestID: %+v", requestID)
	return nil
}

func (s *IdpStorage) RevokeRefreshToken(ctx context.Context, requestID string) error {
	db := infrastructure.Connect()

	result := db.Where("signature=?", requestID).Delete(&models.RefreshToken{})
	if result.Error != nil {
		log.Printf("Error occurred in RevokeRefreshToken: %+v", result.Error)
		return result.Error
	}

	log.Printf("RevokeRefreshToken RequestID: %+v", requestID)
	return nil
}

func (s *IdpStorage) CreateOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) error {
	db := infrastructure.Connect()

	is := models.IDSessionOf(authorizeCode, requester)

	result := db.Create(&is)

	if result.Error != nil {
		log.Printf("Error occurred in CreateOpenIDConnectSession: %+v", result.Error)
		return result.Error
	}

	log.Printf("CreateOpenIDConnectSession Signature: %+v", authorizeCode)
	return nil
}

func (s *IdpStorage) GetOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	db := infrastructure.Connect()

	var is models.IDSession
	result := db.Where("signature=?", authorizeCode).First(&is)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No record found for signature: %s", authorizeCode)
			return nil, fosite.ErrNotFound
		}
		log.Printf("Error occurred in GetOpenIDConnectSession: %+v", result.Error)
		return nil, result.Error
	}

	log.Printf("GetOpenIDConnectSession Signature: %+v", authorizeCode)
	return is.ToRequester(), nil
}

func (s *IdpStorage) DeleteOpenIDConnectSession(_ context.Context, authorizeCode string) error {
	db := infrastructure.Connect()

	result := db.Where("signature=?", authorizeCode).Delete(&models.IDSession{})
	if result.Error != nil {
		log.Printf("Error occurred in DeleteOpenIDConnectSession: %+v", result.Error)
		return result.Error
	}

	log.Printf("DeleteOpenIDConnectSession Signature: %+v", authorizeCode)
	return nil
}
