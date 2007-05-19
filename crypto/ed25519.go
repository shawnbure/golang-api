package crypto

import (
	libed25519 "crypto/ed25519"
)

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
