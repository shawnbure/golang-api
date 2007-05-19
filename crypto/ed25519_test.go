package crypto

import (
	libed25519 "crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSeed(t *testing.T) {
	key := generateSeed()
	t.Log(key)
}

func TestVerifyForCorrectSignature_ShouldPass(t *testing.T) {
	t.Parallel()
	pubKey := "XXnByecIEjQ8Ir/10T/YnCGWX6W48BW+fgmF+PP7iWQ="
	sig := "TJwsQOAbAw9n0twlVXJ2P7FmthrVWaIX5N7j5j6ebxPY0FgpTbRWm7TkbN1jepvQvAXQpsAp8ZLR5OseZnVjBQ=="

	pubKeyBytes, _ := base64.StdEncoding.DecodeString(pubKey)
	sigBytes, _ := base64.StdEncoding.DecodeString(sig)
	msg := []byte("msg")

	err := VerifySignature(pubKeyBytes, msg, sigBytes)
	assert.Nil(t, err)
}

func TestVerifyForIncorrectPubKeyLength_ShouldErr(t *testing.T) {
	t.Parallel()
	pubKey := "XXnByecIEjQ8Ir/10T/YnCGWX6W48BW+fgmF+PP7iWQQ=="
	sig := "TJwsQOAbAw9n0twlVXJ2P7FmthrVWaIX5N7j5j6ebxPY0FgpTbRWm7TkbN1jepvQvAXQpsAp8ZLR5OseZnVjBQ=="

	pubKeyBytes, _ := base64.StdEncoding.DecodeString(pubKey)
	sigBytes, _ := base64.StdEncoding.DecodeString(sig)
	msg := []byte("msg")

	err := VerifySignature(pubKeyBytes, msg, sigBytes)
	assert.Equal(t, ErrInvalidPublicKey, err)
}

func TestVerifyForIncorrectSig_ShouldErr(t *testing.T) {
	t.Parallel()
	pubKey := "XXnByecIEjQ8Ir/10T/YnCGWX6W48BW+fgmF+PP7iWQ="
	sig := "TJwsQOAbAw9n0twlVXJ2P7FmthrVWaIX5n7j5j6ebxPY0FgpTbRWm7TkbN1jepvQvAXQpsAp8ZLR5OseZnVjBQ=="

	pubKeyBytes, _ := base64.StdEncoding.DecodeString(pubKey)
	sigBytes, _ := base64.StdEncoding.DecodeString(sig)
	msg := []byte("msg")

	err := VerifySignature(pubKeyBytes, msg, sigBytes)
	assert.Equal(t, ErrInvalidSignature, err)
}

func TestGibKeySir_KeySignatureWillBeVerified(t *testing.T) {
	t.Parallel()

	seed := "202d2274940909b4f3c23691c857d7d3352a0574cfb96efbf1ef90cbc66e2cbc"
	msg := []byte("all your tokens are belong to us, kind ser")

	seedBytes, _ := hex.DecodeString(seed)

	sk := GibKeySir(seedBytes)
	pk := sk[libed25519.PublicKeySize:]

	sig, _ := SignPayload(sk, msg)

	verifyErr := VerifySignature(pk, msg, sig)
	require.Nil(t, verifyErr)
}
