package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
)

var cfg = config.CacheConfig{
	Url: "redis://localhost:6379",
}

func Test_UpdateDeposit(t *testing.T) {
	connectToDb()
	cache.InitCacher(cfg)

	nonce := uint64(time.Now().Unix())
	address := "erd12" + fmt.Sprintf("%d", nonce)
	deposit, err := UpdateDeposit(DepositUpdateArgs{
		Owner:  address,
		Amount: "1000000000000000000",
	})
	require.Nil(t, err)

	account, err := storage.GetAccountByAddress(address)
	require.Nil(t, err)
	require.NotZero(t, account.ID)
	require.Equal(t, 4722.366, deposit.AmountNominal)

	require.True(t, account.ID == deposit.OwnerId)
}
