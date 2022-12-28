package keys

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const encodeSeparator = "$"

var FailedToDecode = errors.New("failed to decode public key from string")

type PublicKey struct {
	key *rsa.PublicKey
}

func (k *PublicKey) VerifySignature(message, signature []byte) error {
	hashed := sha256.Sum256(message)

	return rsa.VerifyPKCS1v15(k.key, crypto.SHA256, hashed[:], signature)
}

func (k *PublicKey) Encode() string {
	nEncoded := base64.StdEncoding.EncodeToString(k.key.N.Bytes())
	return fmt.Sprintf("%v%v%v", k.key.E, encodeSeparator, nEncoded)
}

func DecodePublicKey(source string) (*PublicKey, error) {
	parts := strings.SplitN(source, encodeSeparator, -1)
	if len(parts) != 2 {
		return nil, FailedToDecode
	}

	e, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	n := new(big.Int)
	nBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	n.SetBytes(nBytes)

	return &PublicKey{
		key: &rsa.PublicKey{
			N: n,
			E: e,
		},
	}, nil
}
