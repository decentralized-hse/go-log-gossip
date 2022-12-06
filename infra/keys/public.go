package keys

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
)

type PublicKey struct {
	key *rsa.PublicKey
}

func (k *PublicKey) VerifySignature(message, signature []byte) error {
	hashed := sha256.Sum256(message)

	return rsa.VerifyPKCS1v15(k.key, crypto.SHA256, hashed[:], signature)
}
