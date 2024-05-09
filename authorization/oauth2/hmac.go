package oauth2

import (
	"context"
	"crypto/sha512"
	"hash"

	"github.com/ory/fosite/token/hmac"
)

type hmacStrategyConfigurator struct {
	Secret []byte
}

func (h *hmacStrategyConfigurator) GetGlobalSecret(_ context.Context) ([]byte, error) {
	return h.Secret, nil
}

func (h *hmacStrategyConfigurator) GetHMACHasher(_ context.Context) func() hash.Hash {
	return sha512.New512_256
}

func (h *hmacStrategyConfigurator) GetRotatedGlobalSecrets(_ context.Context) ([][]byte, error) {
	return nil, nil
}

func (h *hmacStrategyConfigurator) GetTokenEntropy(_ context.Context) int {
	return 32 //nolint:gomnd
}

var _ hmac.HMACStrategyConfigurator = (*hmacStrategyConfigurator)(nil)
