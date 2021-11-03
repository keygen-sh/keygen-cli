package ed25519ph

import (
	"crypto"

	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
)

const (
	SigningKeySize = ed25519.PrivateKeySize
	VerifyKeySize  = ed25519.PublicKeySize
	SignatureSize  = ed25519.SignatureSize
)

type SigningKey = ed25519.PrivateKey
type VerifyKey = ed25519.PublicKey

func GenerateKey() (VerifyKey, SigningKey, error) {
	return ed25519.GenerateKey(nil)
}

func Sign(privateKey SigningKey, digest []byte) ([]byte, error) {
	return privateKey.Sign(nil, digest, &ed25519.Options{Hash: crypto.SHA512})
}

func Verify(publicKey VerifyKey, digest []byte, sig []byte) bool {
	return ed25519.VerifyWithOptions(publicKey, digest, sig, &ed25519.Options{Hash: crypto.SHA512})
}
