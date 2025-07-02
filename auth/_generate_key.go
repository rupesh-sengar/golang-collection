package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Errorf("failed to generate private key: %v", err))
	}

	// 🔐 Encode private key to PKCS#8 PEM
	privateBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(fmt.Errorf("failed to marshal private key: %v", err))
	}

	privatePem := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateBytes,
	}

	err = os.WriteFile("private_oci.pem", pem.EncodeToMemory(privatePem), 0600)
	if err != nil {
		panic(fmt.Errorf("failed to write private.pem: %v", err))
	}

	// 🔓 Encode public key to PKIX PEM (for frontend)
	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(fmt.Errorf("failed to marshal public key: %v", err))
	}

	publicPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	}

	err = os.WriteFile("public_oci.pem", pem.EncodeToMemory(publicPem), 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write public.pem: %v", err))
	}

	fmt.Println("✅ private.pem and public.pem generated successfully")
}
