package services

import (
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_GetOrCreateAccount(t *testing.T) {
	connectToDb()

	account, err := GetOrCreateAccount("address")
	require.Nil(t, err)
	require.Equal(t, account.Address, "address")

	accountRead, err := storage.GetAccountByAddress("address")
	require.Nil(t, err)
	require.Equal(t, accountRead.Address, "address")
}
