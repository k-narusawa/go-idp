package oauth2

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"sync"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	fositeoauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/hmac"
	"github.com/ory/fosite/token/jwt"
)

var (
	config = &fosite.Config{
		IDTokenIssuer:              "http://locahost:3846",
		SendDebugMessagesToClients: true,
		ScopeStrategy:              fosite.ExactScopeStrategy,
		RedirectSecureChecker:      fosite.IsRedirectURISecureStrict,
		AllowedPromptValues:        []string{"none"},
		TokenURL:                   "http://locahost:3846/oauth2/token",
		AccessTokenLifespan:        time.Minute * 30,
		AccessTokenIssuer:          "http://locahost:3846",
		RefreshTokenScopes:         []string{"offline"},
		RefreshTokenLifespan:       time.Hour * 24,
		AuthorizeCodeLifespan:      time.Minute * 1,
	}

	secret = []byte("some-cool-secret-that-is-32bytes")

	privateKey, _ = rsa.GenerateKey(rand.Reader, 2048)

	getPrivateKey = func(context.Context) (interface{}, error) {
		return privateKey, nil
	}

	hmacStrategy = &hmac.HMACStrategy{
		Mutex:  sync.Mutex{},
		Config: &hmacStrategyConfigurator{Secret: secret},
	}

	oAuth2HMACStrategy = &fositeoauth2.HMACSHAStrategy{
		Enigma: hmacStrategy,
		Config: config,
	}
)

func NewOauth2Provider() fosite.OAuth2Provider {
	// var jwtStrategy = compose.NewOAuth2JWTStrategy(getPrivateKey, oAuth2HMACStrategy, config)
	return compose.Compose(
		config,
		NewIdpStorage(),
		&compose.CommonStrategy{
			CoreStrategy: oAuth2HMACStrategy,
			// CoreStrategy:               jwtStrategy,
			OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(getPrivateKey, config),
			Signer: &jwt.DefaultSigner{
				GetPrivateKey: getPrivateKey,
			},
		},

		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.OAuth2RefreshTokenGrantFactory,

		compose.OpenIDConnectExplicitFactory,

		compose.OAuth2TokenIntrospectionFactory,
	)
}
