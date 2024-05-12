package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/k-narusawa/go-idp/logger"
)

func GenerateKey(logger *logger.Logger) {
	if isExists("cert/private.pem") && isExists("cert/public.pem") {
		logger.Info("Key files are already exists. Skip generating key.")
		return
	}

	logger.Info("Generating RSA key pair...")

	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)
	publicKey := key.PublicKey

	savePEMKey("cert/private.pem", key)
	savePublicPEMKey("cert/public.pem", publicKey)
}

func isExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
