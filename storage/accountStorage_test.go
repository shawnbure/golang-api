package storage

import (
	"github.com/erdsea/erdsea-api/data"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AddNewAccount(t *testing.T) {
	connectToDb(t)

	account := defaultAccount()
	error := addNewAccount(&account)
	require.Nil(t, error)

	var accountRead data.Collection
	txRead := GetDB().Last(&accountRead)

	require.Nil(t, txRead.Error)
	assert.Equal(t, accountRead, account)
}

func Test_GetAccountById(t *testing.T) {
	connectToDb(t)

	account := defaultAccount()
	error := addNewAccount(&account)
	require.Nil(t, error)

	accountRead, err := getAccountById(account.ID)
	require.Nil(t, err)
	assert.Equal(t, accountRead.ID, account.ID)
}

func Test_GetAccountByAddress(t *testing.T) {
	connectToDb(t)

	address := "unique_erd_addr"
	account := defaultAccount()
	account.Address = address
	error := addNewAccount(&account)
	require.Nil(t, error)

	retrievedAccount, error := getAccountByAddress(address)
	require.Nil(t, error)
	require.GreaterOrEqual(t, retrievedAccount.Address, address)
}

func defaultAccount() data.Account {
	return data.Account{
		Address: "erd123",
	}
}
