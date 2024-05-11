package oauth2

import (
	"context"
	"crypto/rsa"
	"os"
	"sync"
	"time"

	"github.com/k-narusawa/go-idp/logger"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/token/hmac"
	"github.com/ory/fosite/token/jwt"
	"gopkg.in/yaml.v2"
)

type Oauth2Config struct {
	Issuer                string        `yaml:"issuer"`
	AccessTokenLifespan   time.Duration `yaml:"access_token_lifespan"`
	RefreshTokenLifespan  time.Duration `yaml:"refresh_token_lifespan"`
	AuthorizeCodeLifespan time.Duration `yaml:"authorize_code_lifespan"`
	Hmac                  struct {
		Secret string `yaml:"secret"`
	} `yaml:"hmac"`
}

func NewOauth2Provider(
	privateKey *rsa.PrivateKey,
	logger logger.Logger,
) fosite.OAuth2Provider {
	content, err := os.ReadFile("authorization/oauth2/config.yml")
	if err != nil {
		panic(err)
	}

	var oc Oauth2Config
	if err = yaml.Unmarshal(content, &oc); err != nil {
		panic(err)
	}

	var (
		secret = []byte(oc.Hmac.Secret)

		getPrivateKey = func(context.Context) (interface{}, error) {
			return privateKey, nil
		}

		hmacStrategy = &hmac.HMACStrategy{
			Mutex:  sync.Mutex{},
			Config: &hmacStrategyConfigurator{Secret: secret},
		}

		config = &fosite.Config{
			IDTokenIssuer:              oc.Issuer,
			SendDebugMessagesToClients: true,
			ScopeStrategy:              fosite.ExactScopeStrategy,
			RedirectSecureChecker:      fosite.IsRedirectURISecureStrict,
			AllowedPromptValues:        []string{"none", "login"},
			TokenURL:                   "http://locahost:3846/oauth2/token",
			AccessTokenLifespan:        time.Duration(oc.AccessTokenLifespan.Seconds()),
			AccessTokenIssuer:          oc.Issuer,
			RefreshTokenScopes:         []string{"offline"},
			RefreshTokenLifespan:       time.Duration(oc.RefreshTokenLifespan.Seconds()),
			AuthorizeCodeLifespan:      time.Duration(oc.AuthorizeCodeLifespan.Seconds()),
			GlobalSecret:               secret,
		}

		oAuth2HMACStrategy = &HMACSHAStrategy{
			Enigma: hmacStrategy,
			Config: config,
		}
	)

	// var jwtStrategy = compose.NewOAuth2JWTStrategy(getPrivateKey, oAuth2HMACStrategy, config)
	return compose.Compose(
		config,
		NewIdpStorage(logger),
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
		compose.OAuth2TokenRevocationFactory,

		compose.OAuth2PKCEFactory,
	)
}
