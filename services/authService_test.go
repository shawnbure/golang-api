package services

import (
	libed25519 "crypto/ed25519"
	"encoding/hex"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CreateAndRefreshBeforeExpireShouldNotWork(t *testing.T) {
	seed := "202d2274940909b4f3c23691c857d7d3352a0574cfb96efbf1ef90cbc66e2cbc"
	msg := []byte("all your tokens are belong to us, kind ser")

	seedBytes, _ := hex.DecodeString(seed)

	sk := crypto.NewEdKey(seedBytes)
	pk := sk[libed25519.PublicKeySize:]

	sig, _ := crypto.SignPayload(sk, msg)
	verifyErr := crypto.VerifySignature(pk, msg, sig)
	require.Nil(t, verifyErr)

	service, err := NewAuthService(config.AuthConfig{
		JwtSecret:     "bitcoin-to-1-milly",
		JwtIssuer:     "localhost:5000",
		JwtKeySeedHex: "d6592724167553acf9c8cba9a7dbc7f514efc757d7906546cecfdfc5d4c2e8d1",
		JwtExpiryMins: 15,
	})
	require.Nil(t, err)

	jwt, refresh, err := service.CreateToken(pk, sig, msg)
	require.Nil(t, err)

	// Should err because the token is still valid.
	jwt, refresh, err = service.RefreshToken(jwt, refresh)
	require.NotNil(t, err)

	jwt, refresh, err = service.RefreshToken(jwt, refresh)
	require.NotNil(t, err)

	jwt, refresh, err = service.RefreshToken(jwt, refresh)
	require.NotNil(t, err)

	jwt, refresh, err = service.RefreshToken(jwt, refresh)
	require.NotNil(t, err)
}

func Test_CreateAndRefreshAfterExpireShouldWork(t *testing.T) {
	seed := "202d2274940909b4f3c23691c857d7d3352a0574cfb96efbf1ef90cbc66e2cbc"
	msg := []byte("all your tokens are belong to us, kind ser")

	seedBytes, _ := hex.DecodeString(seed)

	sk := crypto.NewEdKey(seedBytes)
	pk := sk[libed25519.PublicKeySize:]

	sig, _ := crypto.SignPayload(sk, msg)
	verifyErr := crypto.VerifySignature(pk, msg, sig)
	require.Nil(t, verifyErr)

	service, err := NewAuthService(config.AuthConfig{
		JwtSecret:     "bitcoin-to-1-milly",
		JwtIssuer:     "localhost:5000",
		JwtKeySeedHex: "d6592724167553acf9c8cba9a7dbc7f514efc757d7906546cecfdfc5d4c2e8d1",
		JwtExpiryMins: -1,
	})
	require.Nil(t, err)

	jwt, refresh, err := service.CreateToken(pk, sig, msg)
	require.Nil(t, err)

	// Should succeed because token expired.
	jwt, refresh, err = service.RefreshToken(jwt, refresh)
	require.Nil(t, err)

	jwt, refresh, err = service.RefreshToken(jwt, refresh)
	require.Nil(t, err)

	jwt, refresh, err = service.RefreshToken(jwt, refresh)
	require.Nil(t, err)

	jwt, refresh, err = service.RefreshToken(jwt, refresh)
	require.Nil(t, err)
}
