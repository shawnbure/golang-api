package crypto

import (
	libed25519 "crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"testing"

	erdgoData "github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNewEdKey_KeySignatureWillBeVerified(t *testing.T) {
	t.Parallel()

	seed := "202d2274940909b4f3c23691c857d7d3352a0574cfb96efbf1ef90cbc66e2cbc"
	msg := []byte("some test message")

	seedBytes, _ := hex.DecodeString(seed)

	sk := NewEdKey(seedBytes)
	pk := sk[libed25519.PublicKeySize:]

	sig, _ := SignPayload(sk, msg)

	verifyErr := VerifySignature(pk, msg, sig)
	require.Nil(t, verifyErr)
}

func Test_ElrondGoCopyPasted(t *testing.T) {
	address, err := erdgoData.NewAddressFromBech32String("erd19pht2w242wcj0x9gq3us86dtjrrfe3wk8ffh5nhdemf0mce6hsmsupxzlq")
	require.Nil(t, err)

	message := []byte("test message")
	sig, err := hex.DecodeString("ec7a27cb4b23641ae62e3ea96d5858c8142e20d79a6e1710037d1c27b0d138d7452a98da93c036b2b47ee587d4cb4af6ae24c358f3f5f74f85580f45e072280b")
	require.Nil(t, err)

	erdMsg := ComputeElrondSignableMessage(message)
	err = VerifySignature(address.AddressBytes(), erdMsg, sig)
	require.Nil(t, err)
}

func Test_WebWalletRouteSignature(t *testing.T) {
	address, err := erdgoData.NewAddressFromBech32String("erd17s2pz8qrds6ake3qwheezgy48wzf7dr5nhdpuu2h4rr4mt5rt9ussj7xzh")
	require.Nil(t, err)

	message, err := hex.DecodeString("af8ffd30add45b0b7299497e41b3599c5acf81ce2e5989751950f4c25ec94581")
	require.Nil(t, err)

	sig, err := hex.DecodeString("96cb38a3b85fa0adcf2bba88c2453907323faceb2225869d99a42e1c9f65a8d822cb607723b23e5c62910122601ba0094ba043eec7a0c89ba4045e357fbee107")
	require.Nil(t, err)

	err = VerifySignature(address.AddressBytes(), message, sig)
	require.Nil(t, err)
}
