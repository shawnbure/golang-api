package crypto

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	secret = "bitcoin-to-1-milly"
	issuer = "api.erdsea.io"

	addr = "erd111"
)

func TestGenerateJwt_ForTestExpiry25(t *testing.T) {
	t.Parallel()

	jwt, _ := GenerateJwt(addr, secret, issuer, 25)
	t.Log(jwt)
}

func TestGenerateJwt_ShouldValidateThenParseClaims(t *testing.T) {
	t.Parallel()

	jwt, err := GenerateJwt(addr, secret, issuer, 5)
	require.Nil(t, err)

	fmt.Println(jwt)

	c, err := ValidateJwt(jwt, secret)
	require.Nil(t, err)
	require.Equal(t, addr, c.Address)
	require.Equal(t, issuer, c.Issuer)
}

func TestGenerateJwt_ShouldNotValidateWrongExpiry(t *testing.T) {
	t.Parallel()

	jwt, err := GenerateJwt(addr, secret, issuer, 0)
	require.Nil(t, err)

	time.Sleep(time.Second * 1)

	c, err := ValidateJwt(jwt, secret)
	require.NotNil(t, err)
	require.True(t, c.Address == "")
}

func TestGenerateJwt_ShouldNotValidateWrongSecret(t *testing.T) {
	t.Parallel()

	jwt, err := GenerateJwt(addr, "bad", issuer, 5)
	require.Nil(t, err)

	c, err := ValidateJwt(jwt, secret)
	require.NotNil(t, err)
	require.True(t, c.Address == "")
}

func TestGetClaims_ShouldReturnForValidToken(t *testing.T) {
	t.Parallel()

	jwt, err := GenerateJwt(addr, secret, issuer, 5)
	require.Nil(t, err)

	c, err := GetClaims(jwt, secret, true)
	require.Nil(t, err)
	require.Equal(t, c.Address, addr)
}

func TestGetClaims_ShouldReturnForExpiredToken(t *testing.T) {
	t.Parallel()

	jwt, err := GenerateJwt(addr, secret, issuer, 0)
	require.Nil(t, err)

	time.Sleep(time.Second * 1)

	c, err := GetClaims(jwt, secret, false)
	require.Nil(t, err)
	require.Equal(t, c.Address, addr)
}
