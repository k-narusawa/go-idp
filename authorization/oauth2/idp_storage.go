package oauth2

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"time"

	"github.com/k-narusawa/go-idp/adapter/gateways"
	"github.com/k-narusawa/go-idp/domain/models"
	"github.com/k-narusawa/go-idp/logger"

	"github.com/ory/fosite"
	"gorm.io/gorm"
)

type IdpStorage struct {
	logger             logger.Logger
	Clients            []models.Client
	OidcSessions       []models.OidcSession
	AuthorizationCodes []models.AuthorizationCode
	AccessTokens       []models.AccessToken
	RefreshTokens      []models.RefreshToken
	PKCES              []models.PKCE
}

func NewIdpStorage(logger logger.Logger) *IdpStorage {
	is := IdpStorage{
		logger: logger,
	}
	return &is
}

func (s *IdpStorage) CreateClient(_ context.Context, client fosite.Client) {
	s.logger.Info("CreateClient",
		slog.String("client_id", client.GetID()),
	)

	db := gateways.Connect()

	c := models.ClientOf(client)
	result := db.Create(&c)
	if result.Error != nil {
		s.logger.Error("Error occurred in CreateClient", result.Error)
	}
}

func (s *IdpStorage) GetClient(_ context.Context, id string) (fosite.Client, error) {
	s.logger.Info("GetClient",
		slog.String("id", id),
	)

	db := gateways.Connect()

	var c models.Client
	res := db.Where("id=?", id).First(&c)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			s.logger.Warn("No record found for id",
				slog.String("id", id),
			)
			return nil, fosite.ErrNotFound
		}
		s.logger.Error("Error occurred in GetClient", res.Error)
		return nil, res.Error
	}

	return models.CastToClient(c), nil
}

func (s *IdpStorage) ClientAssertionJWTValid(_ context.Context, jti string) error {
	s.logger.Info("ClientAssertionJWTValid",
		slog.String("jti", jti),
	)
	return nil
}

func (s *IdpStorage) SetClientAssertionJWT(_ context.Context, jti string, exp time.Time) error {
	s.logger.Info("SetClientAssertionJWT",
		slog.String("jti", jti),
	)
	return nil
}

func (s *IdpStorage) CreateOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) error {
	s.logger.Info("CreateOpenIDConnectSession",
		slog.String("signature", authorizeCode),
	)

	db := gateways.Connect()

	is := models.IDSessionOf(authorizeCode, requester)

	result := db.Create(&is)

	if result.Error != nil {
		s.logger.Error("Error occurred in CreateOpenIDConnectSession",
			result.Error,
		)
		return result.Error
	}

	return nil
}

func (s *IdpStorage) GetOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	s.logger.Info("GetOpenIDConnectSession",
		slog.String("signature", authorizeCode),
	)

	db := gateways.Connect()

	var is models.OidcSession
	result := db.Where("signature=?", authorizeCode).First(&is)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			s.logger.Warn("No record found for signature",
				slog.String("signature", authorizeCode),
			)
			return nil, fosite.ErrNotFound
		}
		s.logger.Error("Error occurred in GetOpenIDConnectSession",
			result.Error,
		)
		return nil, result.Error
	}

	return is.ToRequester(), nil
}

func (s *IdpStorage) DeleteOpenIDConnectSession(_ context.Context, authorizeCode string) error {
	// Deprecated
	s.logger.Info("DeleteOpenIDConnectSession",
		slog.String("signature", authorizeCode),
		slog.String("operation", "DEPRECATED"),
	)
	return nil
}

func (s *IdpStorage) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	s.logger.Info("CreateAccessTokenSession",
		slog.String("signature", signature),
	)
	db := gateways.Connect()

	at := models.AccessTokenOf(signature, request)

	result := db.Create(&at)

	if result.Error != nil {
		s.logger.Error("Error occurred in CreateAccessTokenSession",
			result.Error,
		)
		return result.Error
	}

	return nil
}

func (s *IdpStorage) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	s.logger.Info("GetAccessTokenSession",
		slog.String("signature", signature),
	)

	db := gateways.Connect()

	var at models.AccessToken

	result := db.
		Preload("Client").
		Where("signature=?", signature).
		Find(&at)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			s.logger.Warn("No record found for signature",
				slog.String("signature", signature),
			)
			return nil, fosite.ErrNotFound
		}
		s.logger.Error("Error occurred in GetAccessTokenSession",
			result.Error,
		)
		return nil, result.Error
	}

	return at.ToRequester(), nil
}

func (s *IdpStorage) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	s.logger.Info("DeleteAccessTokenSession",
		slog.String("signature", signature),
	)

	db := gateways.Connect()

	result := db.Where("signature=?", signature).Delete(&models.AccessToken{})
	if result.Error != nil {
		s.logger.Error("Error occurred in DeleteAccessTokenSession",
			result.Error,
		)
		return result.Error
	}

	return nil
}

func (s *IdpStorage) CreateAuthorizeCodeSession(_ context.Context, code string, req fosite.Requester) error {
	s.logger.Info("CreateAuthorizeCodeSession",
		slog.String("code", code),
	)

	db := gateways.Connect()

	ac := models.AuthorizationCodeOf(code, req)

	result := db.Create(&ac)
	if result.Error != nil {
		s.logger.Error("Error occurred in CreateAuthorizeCodeSession",
			result.Error,
		)
		return result.Error
	}

	return nil
}

func (s *IdpStorage) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {
	s.logger.Info("GetAuthorizeCodeSession",
		slog.String("code", code),
	)
	db := gateways.Connect()

	var ac models.AuthorizationCode
	ar := db.
		Preload("Client").
		Where("signature=?", code).
		Find(&ac)

	if ar.Error != nil {
		s.logger.Error("Error occurred in GetAuthorizeCodeSession",
			ar.Error,
		)
		return nil, ar.Error
	}

	return ac.ToRequester(), nil
}

func (s *IdpStorage) InvalidateAuthorizeCodeSession(ctx context.Context, code string) (err error) {
	s.logger.Info("InvalidateAuthorizeCodeSession",
		slog.String("code", code),
	)

	db := gateways.Connect()

	result := db.Where("signature=?", code).Delete(&models.AuthorizationCode{})
	if result.Error != nil {
		s.logger.Error("Error occurred in InvalidateAuthorizeCodeSession", result.Error)
		return result.Error
	}

	return nil
}

func (s *IdpStorage) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	s.logger.Info("CreateRefreshTokenSession",
		slog.String("signature", signature),
	)
	db := gateways.Connect()

	rt := models.RefreshTokenOf(signature, request)
	result := db.Create(&rt)

	if result.Error != nil {
		s.logger.Error("Error occurred in CreateRefreshTokenSession", result.Error)
		return result.Error
	}

	return nil
}

func (s *IdpStorage) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	s.logger.Info("GetRefreshTokenSession",
		slog.String("signature", signature),
	)

	db := gateways.Connect()

	var rt models.RefreshToken
	result := db.
		Preload("Client").
		Where("signature=?", signature).First(&rt)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			s.logger.Warn("No record found for signature",
				slog.String("signature", signature),
			)
			return nil, fosite.ErrNotFound
		}
		s.logger.Error("Error occurred in GetRefreshTokenSession", result.Error)
		return nil, result.Error
	}

	return rt.ToRequester(), nil
}

func (s *IdpStorage) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	s.logger.Info("DeleteRefreshTokenSession",
		slog.String("signature", signature),
	)

	db := gateways.Connect()

	result := db.Where("signature=?", signature).Delete(&models.RefreshToken{})
	if result.Error != nil {
		s.logger.Error("Error occurred in DeleteRefreshTokenSession", result.Error)
		return result.Error
	}

	return nil
}

func (s *IdpStorage) RevokeAccessToken(ctx context.Context, requestID string) error {
	s.logger.Info("RevokeAccessToken",
		slog.String("requestID", requestID),
	)

	db := gateways.Connect()

	result := db.Where("request_id=?", requestID).Delete(&models.AccessToken{})
	if result.Error != nil {
		log.Printf("Error occurred in RevokeAccessToken: %+v", result.Error)
		return result.Error
	}

	return nil
}

func (s *IdpStorage) RevokeRefreshToken(ctx context.Context, requestID string) error {
	s.logger.Info("RevokeRefreshToken",
		slog.String("requestID", requestID),
	)

	db := gateways.Connect()

	result := db.Where("request_id=?", requestID).Delete(&models.RefreshToken{})
	if result.Error != nil {
		s.logger.Error("Error occurred in RevokeRefreshToken", result.Error)
		return result.Error
	}

	return nil
}

func (s *IdpStorage) RevokeRefreshTokenMaybeGracePeriod(ctx context.Context, requestID string, signature string) error {
	// no configuration option is available; grace period is not available with memory store
	return s.RevokeRefreshToken(ctx, requestID)
}

func (s *IdpStorage) CreatePKCERequestSession(_ context.Context, code string, req fosite.Requester) error {
	s.logger.Info("CreatePKCERequestSession",
		slog.String("signature", code),
	)

	db := gateways.Connect()

	pkce := models.PKCEOf(code, req)

	result := db.Create(&pkce)
	if result.Error != nil {
		s.logger.Error("Error occurred in CreatePKCERequestSession", result.Error)
		return result.Error
	}

	return nil
}

func (s *IdpStorage) GetPKCERequestSession(_ context.Context, code string, _ fosite.Session) (fosite.Requester, error) {
	s.logger.Info("GetPKCERequestSession",
		slog.String("signature", code),
	)

	db := gateways.Connect()

	var pkce models.PKCE
	result := db.Where("signature=?", code).First(&pkce)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			s.logger.Warn("No record found for signature",
				slog.String("signature", code),
			)
			return nil, fosite.ErrNotFound
		}
		s.logger.Error("Error occurred in GetPKCERequestSession", result.Error)
		return nil, result.Error
	}

	return pkce.ToRequester(), nil
}

func (s *IdpStorage) DeletePKCERequestSession(_ context.Context, code string) error {
	s.logger.Info("DeletePKCERequestSession",
		slog.String("signature", code),
	)

	db := gateways.Connect()

	result := db.Where("signature=?", code).Delete(&models.PKCE{})
	if result.Error != nil {
		s.logger.Error("Error occurred in DeletePKCERequestSession", result.Error)
		return result.Error
	}

	return nil
}
