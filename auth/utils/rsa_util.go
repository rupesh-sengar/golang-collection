package utils

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"github.com/rupesh-sengar/auth/config"
	"strings"
)

func DecryptEncPassword(enc string) (string, error) {
	parts := strings.Split(enc, ":")
	if len(parts) < 4 {
		return "", errors.New("invalid enc_password format")
	}

	encryptedBase64 := parts[3]
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := rsa.DecryptPKCS1v15(nil, config.PrivateKey, encryptedBytes)
	if err != nil {
		return "", err
	}

	payload := string(decryptedBytes)
	split := strings.SplitN(payload, ":", 2)
	if len(split) != 2 {
		return "", errors.New("decrypted payload malformed")
	}

	return split[1], nil
}
