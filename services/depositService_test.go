package services

import (
	"testing"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/config"
	"github.com/stretchr/testify/require"
)

var cacheCfg = config.CacheConfig{
	Url: "redis://localhost:6379",
}

func Test_UpdateDeposit(t *testing.T) {
	connectToDb()
	cache.InitCacher(cacheCfg)

	err := UpdateDeposit(DepositUpdateArgs{
		Owner: "erd1",
		Amount: "100000000000000000000",
	})
	require.Nil(t, err)

	deposit, err := GetDeposit(blockchainCfg.MarketplaceAddress, "erd1")
	require.Nil(t, err)
	require.Equal(t, 1208925.819, deposit)
}
