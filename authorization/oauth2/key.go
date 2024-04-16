package oauth2

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func ReadPrivatekey() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile("cert/private.pem")
	if err != nil {
		return nil, err
	}
	keyblock, _ := pem.Decode(data)
	if keyblock == nil {
		return nil, fmt.Errorf("invalid private key data")
	}

	if keyblock.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key type : %s", keyblock.Type)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyblock.Bytes)

	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
