package services

import (
	"testing"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
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

func Test_SearchAccount(T *testing.T) {
	connectToDb()
	cache.InitCacher(config.CacheConfig{Url: "redis://localhost:6379"})

	acc := &data.Account{
		Name: "this name is uniquee",
	}

	acc.ID = 0
	err := storage.AddAccount(acc)
	require.Nil(T, err)

	acc.ID = 0
	err = storage.AddAccount(acc)
	require.Nil(T, err)

	acc.ID = 0
	err = storage.AddAccount(acc)
	require.Nil(T, err)

	acc.ID = 0
	err = storage.AddAccount(acc)
	require.Nil(T, err)

	acc.ID = 0
	err = storage.AddAccount(acc)
	require.Nil(T, err)

	acc.ID = 0
	err = storage.AddAccount(acc)
	require.Nil(T, err)

	accs, err := GetAccountsWithNameAlike("uniquee", 5)
	require.Nil(T, err)
	require.Equal(T, len(accs), 5)
	require.Equal(T, accs[0].Name, "this name is uniquee")
	require.Equal(T, accs[1].Name, "this name is uniquee")
	require.Equal(T, accs[2].Name, "this name is uniquee")
	require.Equal(T, accs[3].Name, "this name is uniquee")
	require.Equal(T, accs[4].Name, "this name is uniquee")

	accs, err = GetAccountsWithNameAlike("uniquee", 5)
	require.Nil(T, err)
	require.Equal(T, len(accs), 5)
	require.Equal(T, accs[0].Name, "this name is uniquee")
	require.Equal(T, accs[1].Name, "this name is uniquee")
	require.Equal(T, accs[2].Name, "this name is uniquee")
	require.Equal(T, accs[3].Name, "this name is uniquee")
	require.Equal(T, accs[4].Name, "this name is uniquee")

	hits := cache.GetCacher().GetStats().Hits
	require.GreaterOrEqual(T, hits.Load(), int64(1))
}
