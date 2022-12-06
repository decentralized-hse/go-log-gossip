package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
)

type PublicPrivateKeyPair struct {
	key *rsa.PrivateKey
}

func GenerateNewPair() (pair *PublicPrivateKeyPair) {
	key, err := rsa.GenerateKey(rand.Reader, KeyBitSize)
	if err != nil {
		panic(err)
	}

	return &PublicPrivateKeyPair{key}
}

func (p *PublicPrivateKeyPair) SaveToFiles(path string) {
	// https://stackoverflow.com/questions/64104586/use-golang-to-get-rsa-key-the-same-way-openssl-genrsa
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(p.key),
		},
	)

	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(p.key.Public().(*rsa.PublicKey)),
		},
	)

	fullPrivatePath := filepath.Join(path, PrivateKeyFilename)
	if err := os.WriteFile(fullPrivatePath, keyPEM, 0400); err != nil {
		panic(err)
	}

	fullPublicPath := filepath.Join(path, PublicKeyFilename)
	if err := os.WriteFile(fullPublicPath, pubPEM, 0444); err != nil {
		panic(err)
	}
}

func LoadFromFiles(path string) (*PublicPrivateKeyPair, error) {
	fullPrivatePath := filepath.Join(path, PrivateKeyFilename)
	fullPublicPath := filepath.Join(path, PublicKeyFilename)

	privateKey, err := parsePrivateKey(fullPrivatePath)
	if err != nil {
		return nil, err
	}

	publicKey, err := parsePublicKey(fullPublicPath)
	if err != nil {
		return nil, err
	}

	privateKey.PublicKey = *publicKey
	privateKey.Precompute()

	return &PublicPrivateKeyPair{privateKey}, nil
}

func (p *PublicPrivateKeyPair) GetPrivateKey() *PrivateKey {
	return &PrivateKey{p.key}
}

func (p *PublicPrivateKeyPair) GetPublicKey() *PublicKey {
	return &PublicKey{p.key.Public().(*rsa.PublicKey)}
}

func parsePrivateKey(fullPathToPrivateKey string) (privateKey *rsa.PrivateKey, err error) {
	privateKeyRawData, err := os.ReadFile(fullPathToPrivateKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyRawData)
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	return
}

func parsePublicKey(fullPathToPublicKey string) (publicKey *rsa.PublicKey, err error) {
	publicKeyRawData, err := os.ReadFile(fullPathToPublicKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyRawData)
	publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	return
}
