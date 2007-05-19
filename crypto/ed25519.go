package crypto

import (
	"crypto/rand"
	"encoding/hex"

	libed25519 "crypto/ed25519"
)

func generateSeed() string {
	_, sk, _ := libed25519.GenerateKey(rand.Reader)

	return hex.EncodeToString(sk.Seed())
}

func NewEdKey(seed []byte) libed25519.PrivateKey {
	return libed25519.NewKeyFromSeed(seed)
}

func VerifySignature(publicKey, message, signature []byte) error {
	if len(publicKey) != libed25519.PublicKeySize {
		return ErrInvalidPublicKey
	}

	isValidSig := libed25519.Verify(publicKey, message, signature)
	if !isValidSig {
		return ErrInvalidSignature
	}

	return nil
}

func SignPayload(privKey, message []byte) ([]byte, error) {
	if len(privKey) != libed25519.PrivateKeySize {
		return nil, ErrInvalidPrivateKey
	}

	sig := libed25519.Sign(privKey, message)

	return sig, nil
}
