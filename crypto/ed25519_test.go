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
	msg := []byte("all your tokens are belong to us, kind ser")

	seedBytes, _ := hex.DecodeString(seed)

	sk := NewEdKey(seedBytes)
	pk := sk[libed25519.PublicKeySize:]

	sig, _ := SignPayload(sk, msg)

	verifyErr := VerifySignature(pk, msg, sig)
	require.Nil(t, verifyErr)
}

func Test_VerifyDevnetWalletGeneratedSignature(t *testing.T) {
	address, err := erdgoData.NewAddressFromBech32String("erd17s2pz8qrds6ake3qwheezgy48wzf7dr5nhdpuu2h4rr4mt5rt9ussj7xzh")
	require.Nil(t, err)

	message := []byte("cevaceva")
	sig, err := hex.DecodeString("8722fc7a40c84ab784d7cca3c94a334bd2da82fd55c827e242fe4bc3a7062342d7f61ac037bee380dac1237ea369bc390882059abb965ab98855139dc7745e0c")
	require.Nil(t, err)

	erdMsg := ComputeElrondSignableMessage(message)
	err = VerifySignature(address.AddressBytes(), erdMsg, sig)
	require.Nil(t, err)
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
