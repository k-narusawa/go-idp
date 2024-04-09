package oauth2

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"idp/infrastructure"
	"idp/models"

	"github.com/ory/fosite"
	"gorm.io/gorm"
)

type IdpStorage struct {
	Clients            []models.Client
	Users              []models.User
	IDSessions         sync.Map
	AuthorizationCodes sync.Map
	AccessTokens       sync.Map
	RefreshTokens      sync.Map
}

func NewIdpStorage() *IdpStorage {
	is := IdpStorage{}
	// is.CreateClient(context.TODO(), defaultClient)
	return &is
}

func (s *IdpStorage) CreateClient(_ context.Context, client fosite.Client) {
	db := infrastructure.Connect()
	db.Create(&models.Client{
		ID:             "my-client",
		Secret:         []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
		RotatedSecrets: `$2y$10$X51gLxUQJ.hGw1epgHTE5u0bt64xM0COU7K9iAp.OFg8p2pUd.1zC `,
		RedirectURIs:   "http://localhost:3846/callback",
		ResponseTypes:  "id_token,code,token,id_token token,code id_token,code token,code id_token token",
		GrantTypes:     "implicit,refresh_token,authorization_code,password,client_credentials",
		Scopes:         "fosite,openid,offline",
	})
}

func (s *IdpStorage) GetClient(_ context.Context, id string) (fosite.Client, error) {
	db := infrastructure.Connect()
	log.Printf("client_id: %+v", id)

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
	s.AccessTokens.Store(signature, request)
	return nil
}

func (s *IdpStorage) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	at, ok := s.AccessTokens.Load(signature)
	if ok {
		return at.(fosite.Requester), nil
	}

	return nil, fosite.ErrNotFound
}

func (s *IdpStorage) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	s.AccessTokens.Delete(signature)
	return nil
}

func (s *IdpStorage) CreateAuthorizeCodeSession(_ context.Context, code string, req fosite.Requester) error {
	log.Printf("CreateAuthorizeCodeSession: %+v", code)
	s.AuthorizationCodes.Store(code, req)

	return nil
}

func (s *IdpStorage) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {
	log.Printf("CreateAuthorizeCodeSession: %+v", code)
	ac, ok := s.AuthorizationCodes.Load(code)

	if ok {
		return ac.(fosite.Requester), nil
	}

	return nil, fosite.ErrNotFound
}

func (s *IdpStorage) InvalidateAuthorizeCodeSession(ctx context.Context, code string) (err error) {
	log.Printf("CreateAuthorizeCodeSession: %+v", code)
	return nil
}

func (s *IdpStorage) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	s.RefreshTokens.Store(signature, request)
	return nil
}

func (s *IdpStorage) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	at, ok := s.RefreshTokens.Load(signature)
	if ok {
		return at.(fosite.Requester), nil
	}

	return nil, fosite.ErrNotFound
}

func (s *IdpStorage) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	s.RefreshTokens.Delete(signature)
	return nil
}

func (s *IdpStorage) RevokeAccessToken(ctx context.Context, requestID string) error {
	return nil
}

func (s *IdpStorage) RevokeRefreshToken(ctx context.Context, requestID string) error {
	return nil
}

func (s *IdpStorage) CreateOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) error {
	s.IDSessions.Store(authorizeCode, requester)
	return nil
}

func (s *IdpStorage) GetOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	ac, ok := s.IDSessions.Load(authorizeCode)
	if ok {
		return ac.(fosite.Requester), nil
	}
	return nil, fosite.ErrNotFound
}

func (s *IdpStorage) DeleteOpenIDConnectSession(_ context.Context, authorizeCode string) error {
	return nil
}
