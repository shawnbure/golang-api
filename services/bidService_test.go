package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/stretchr/testify/require"
)

func Test_PlaceBid(t *testing.T) {
	connectToDb()
	cache.InitCacher(cacheCfg)

	nonce := uint64(time.Now().Unix())
	token := entities.Token{
		TokenID: "TEST",
		Nonce:   nonce,
	}
	err := storage.AddToken(&token)
	require.Nil(t, err)

	address := "erd12" + fmt.Sprintf("%d", nonce)
	err = UpdateDeposit(DepositUpdateArgs{
		Owner:  address,
		Amount: "1000000000000000000",
	})
	require.Nil(t, err)

	offer, err := PlaceBid(PlaceBidArgs{
		Offeror: address,
		TokenId: "TEST",
		Amount:  "1000000000000000000",
		Nonce:   nonce,
	})
	require.Nil(t, err)
	require.Equal(t, token.ID, offer.TokenID)
}
