package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
)

var PrivateKey *rsa.PrivateKey

func LoadRSAKeys() error {
	ENCODED_PRIVATE_KEY := os.Getenv("PRIVATE_KEY")

	PRIVATE_KEY, _ := base64.StdEncoding.DecodeString(ENCODED_PRIVATE_KEY)
	block, _ := pem.Decode(PRIVATE_KEY)
	if block == nil || block.Type != "PRIVATE KEY" {
		return errors.New("failed to parse PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	if err != nil {
		return err
	}

	PrivateKey = key.(*rsa.PrivateKey)
	return nil
}
