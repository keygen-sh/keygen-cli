package ed25519ph

import (
	"crypto"

	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
)

type PrivateKey = ed25519.PrivateKey
type PublicKey = ed25519.PublicKey

func GenerateKey() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(nil)
}

func Sign(privateKey ed25519.PrivateKey, digest []byte) ([]byte, error) {
	return privateKey.Sign(nil, digest, &ed25519.Options{Hash: crypto.SHA512})
}

func Verify(publicKey ed25519.PublicKey, digest []byte, sig []byte) bool {
	return ed25519.VerifyWithOptions(publicKey, digest, sig, &ed25519.Options{Hash: crypto.SHA512})
}
