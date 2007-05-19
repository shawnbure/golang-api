package services

import (
	"github.com/erdsea/erdsea-api/data"
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

func Test_AddOrUpdate(t *testing.T) {
	connectToDb()

	account := data.Account{};
	account.Address = "unique_erd_addr_4real"
	account.Description = "old"
	err := AddOrUpdateAccount(&account)
	require.Nil(t, err)

	retrievedAccount, err := storage.GetAccountByAddress(account.Address)
	require.Nil(t, err)
	require.GreaterOrEqual(t, retrievedAccount.Address, account.Address)
	require.Equal(t, retrievedAccount.Description, "old")

	account.Description = "new"
	err = AddOrUpdateAccount(&account)

	retrievedAccount2, err := storage.GetAccountByAddress(account.Address)
	require.Nil(t, err)
	require.GreaterOrEqual(t, retrievedAccount2.Address, account.Address)
	require.Equal(t, retrievedAccount2.Description, "new")
	require.Equal(t, retrievedAccount.ID, retrievedAccount2.ID)
}
