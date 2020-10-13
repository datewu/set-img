package auth

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
)

// Valid ...
func Valid(token string) (bool, error) {
	sig, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return false, err
	}
	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, digstHash[:], sig)
	return valid, nil
}

// NewToken ...
func NewToken() (string, error) {
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, digstHash[:])
	if err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(sig)
	return token, nil
}
