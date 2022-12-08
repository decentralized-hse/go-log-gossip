package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type PrivateKey struct {
	key *rsa.PrivateKey
}

func (k *PrivateKey) SignMessage(message []byte) (signature []byte, err error) {
	hashed := sha256.Sum256(message)

	signature, err = rsa.SignPKCS1v15(rand.Reader, k.key, crypto.SHA256, hashed[:])

	if err != nil {
		return nil, err
	}

	return
}
