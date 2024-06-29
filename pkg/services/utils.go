package services

import (
	"crypto/rsa"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// Reads an RSA private key from the specified file path
// and returns the parsed private key or an error if the operation fails.
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return nil, err
	}

	return parsedKey, nil
}

// Reads an RSA public key from the specified file path
// and returns the parsed public key or an error if the operation fails.
func loadPublicKey(path string) (*rsa.PublicKey, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
