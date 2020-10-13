package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io/ioutil"
	"os"
)

const (
	digest = "DIGEST_FOR_SIGN which is irrelevant, :)"
)

var (
	privateKey *ecdsa.PrivateKey
	digstHash  = sha256.Sum256([]byte(digest))
)

// InitKeys load/generate the ecdsa private key
func InitKeys(fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return generatePrivateKey(fn)
		}
		return err
	}
	f.Close()
	return loadPrivateKey(fn)
}

func generatePrivateKey(fn string) error {
	pri, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	bs, err := encodePrivateKey(pri)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fn, bs, 0600)
	if err != nil {
		return err
	}
	privateKey = pri
	return nil
}

func loadPrivateKey(fn string) error {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	pri, err := decodePrivateKey(bs)
	if err != nil {
		return err
	}
	privateKey = pri
	return nil
}
