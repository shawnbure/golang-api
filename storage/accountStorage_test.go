package storage

import (
	"github.com/erdsea/erdsea-api/data/entities"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AddAccount(t *testing.T) {
	connectToTestDb()

	account := defaultAccount()
	err := AddAccount(&account)
	require.Nil(t, err)

	var accountRead entities.Account
	txRead := GetDB().Last(&accountRead)

	require.Nil(t, txRead.Error)
	require.Equal(t, accountRead, account)
}

func Test_GetAccountById(t *testing.T) {
	connectToTestDb()

	account := defaultAccount()
	err := AddAccount(&account)
	require.Nil(t, err)

	accountRead, err := GetAccountById(account.ID)
	require.Nil(t, err)
	require.Equal(t, accountRead, &account)
}

func Test_GetAccountByAddress(t *testing.T) {
	connectToTestDb()

	account := defaultAccount()
	account.Address = "unique_erd_addr"
	err := AddAccount(&account)
	require.Nil(t, err)

	retrievedAccount, err := GetAccountByAddress(account.Address)
	require.Nil(t, err)
	require.GreaterOrEqual(t, retrievedAccount.Address, account.Address)
}

func Test_GetAccountsWithNameAlikeWithLimit(t *testing.T) {
	connectToTestDb()

	account := defaultAccount()
	_ = AddAccount(&account)
	account.ID = 0
	_ = AddAccount(&account)

	retrievedAccounts, err := GetAccountsWithNameAlikeWithLimit("default", 5)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(retrievedAccounts), 2)
	require.Equal(t, retrievedAccounts[0].Name, "default")
	require.Equal(t, retrievedAccounts[1].Name, "default")
}

func defaultAccount() entities.Account {
	return entities.Account{
		Address: "erd123",
		Name:    "default",
	}
}
