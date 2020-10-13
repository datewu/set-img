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
	digest     = "DIGEST_FOR_SIGN which is irrelevant, :)"
	priKeyFile = "private_key_for_sign.pem"
)

var (
	privateKey *ecdsa.PrivateKey
	digstHash  = sha256.Sum256([]byte(digest))
)

func init() {
	f, err := os.Open(priKeyFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			generatePrivateKey()
			return
		}
		panic(err)
	}
	f.Close()
	loadPrivateKey()
}

func generatePrivateKey() {
	pri, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	bs, err := encodePrivateKey(pri)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(priKeyFile, bs, 0600)
	if err != nil {
		panic(err)
	}
	privateKey = pri
}

func loadPrivateKey() {
	bs, err := ioutil.ReadFile(priKeyFile)
	if err != nil {
		panic(err)
	}
	pri, err := decodePrivateKey(bs)
	if err != nil {
		panic(err)
	}
	privateKey = pri
}
