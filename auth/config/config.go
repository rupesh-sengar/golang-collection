package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"errors"
	"os"
)

var PrivateKey *rsa.PrivateKey

func LoadRSAKeys() error {
	encodedPrivateKey := os.Getenv("PRIVATE_KEY")
	if encodedPrivateKey == "" {
		return errors.New("PRIVATE_KEY environment variable not set")
	}

	privateKeyBytes, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to base64 decode private key: %w", err)
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return errors.New("failed to parse PEM block containing the private key")
	}
	if block.Type != "PRIVATE KEY" {
		return fmt.Errorf("unexpected key type: %s", block.Type)
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse PKCS8 private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return errors.New("key is not an RSA private key")
	}

	PrivateKey = rsaKey
	return nil
}
