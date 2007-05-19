package storage

import (
	"github.com/erdsea/erdsea-api/data"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AddNewAccount(t *testing.T) {
	connectToTestDb()

	account := defaultAccount()
	err := AddNewAccount(&account)
	require.Nil(t, err)

	var accountRead data.Account
	txRead := GetDB().Last(&accountRead)

	require.Nil(t, txRead.Error)
	require.Equal(t, accountRead, account)
}

func Test_GetAccountById(t *testing.T) {
	connectToTestDb()

	account := defaultAccount()
	err := AddNewAccount(&account)
	require.Nil(t, err)

	accountRead, err := GetAccountById(account.ID)
	require.Nil(t, err)
	require.Equal(t, accountRead, &account)
}

func Test_GetAccountByAddress(t *testing.T) {
	connectToTestDb()

	account := defaultAccount()
	account.Address = "unique_erd_addr"
	err := AddNewAccount(&account)
	require.Nil(t, err)

	retrievedAccount, err := GetAccountByAddress(account.Address)
	require.Nil(t, err)
	require.GreaterOrEqual(t, retrievedAccount.Address, account.Address)
}

func defaultAccount() data.Account {
	return data.Account{
		Address: "erd123",
	}
}
