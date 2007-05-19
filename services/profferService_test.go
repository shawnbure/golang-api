package services

import (
	"fmt"
	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_MakeOffer(t *testing.T) {
	connectToDb()
	cache.InitCacher(cfg)

	nonce := uint64(time.Now().Unix())
	token := entities.Token{
		TokenID: "TEST",
		Nonce:   nonce,
	}
	err := storage.AddToken(&token)
	require.Nil(t, err)

	address := "erd12" + fmt.Sprintf("%d", nonce)
	deposit, err := UpdateDeposit(DepositUpdateArgs{
		Owner:  address,
		Amount: "1000000000000000000",
	})
	require.Nil(t, err)

	offer, err := MakeOffer(MakeOfferArgs{
		OfferorAddress: address,
		TokenId:        "TEST",
		Amount:         "1000000000000000000",
		Nonce:          nonce,
	})
	require.Nil(t, err)
	require.Equal(t, entities.ProfferType("Offer"), offer.Type)
	require.Equal(t, token.ID, offer.TokenID)
	require.Equal(t, deposit.OwnerId, offer.OfferorID)
}

func Test_MakeOfferAcceptOffer(t *testing.T) {
	connectToDb()
	cache.InitCacher(cfg)

	nonce := uint64(time.Now().Unix())
	token := entities.Token{
		TokenID: "TEST",
		Nonce:   nonce,
	}
	err := storage.AddToken(&token)
	require.Nil(t, err)

	address := "erd12" + fmt.Sprintf("%d", nonce)
	deposit, err := UpdateDeposit(DepositUpdateArgs{
		Owner:  address,
		Amount: "1000000000000000000",
	})
	require.Nil(t, err)

	offer, err := MakeOffer(MakeOfferArgs{
		OfferorAddress: address,
		TokenId:        "TEST",
		Amount:         "1000000000000000000",
		Nonce:          nonce,
	})
	offer, err = MakeOffer(MakeOfferArgs{
		OfferorAddress: address,
		TokenId:        "TEST",
		Amount:         "1000000000000000000",
		Nonce:          nonce,
	})
	require.Nil(t, err)
	require.Equal(t, entities.ProfferType("Offer"), offer.Type)
	require.Equal(t, token.ID, offer.TokenID)
	require.Equal(t, deposit.OwnerId, offer.OfferorID)

	AcceptOffer(AcceptOfferArgs{
		OwnerAddress:   address,
		TokenId:        token.TokenID,
		Nonce:          token.Nonce,
		OfferorAddress: address,
		Amount:         "1000000000000000000",
	})
}

func Test_MakeOfferCancelOffer(t *testing.T) {
	connectToDb()
	cache.InitCacher(cfg)

	nonce := uint64(time.Now().Unix())
	token := entities.Token{
		TokenID: "TEST",
		Nonce:   nonce,
	}
	err := storage.AddToken(&token)
	require.Nil(t, err)

	address := "erd12" + fmt.Sprintf("%d", nonce)
	deposit, err := UpdateDeposit(DepositUpdateArgs{
		Owner:  address,
		Amount: "1000000000000000000",
	})
	require.Nil(t, err)

	offer, err := MakeOffer(MakeOfferArgs{
		OfferorAddress: address,
		TokenId:        "TEST",
		Amount:         "1000000000000000000",
		Nonce:          nonce,
	})
	require.Nil(t, err)
	require.Equal(t, entities.ProfferType("Offer"), offer.Type)
	require.Equal(t, token.ID, offer.TokenID)
	require.Equal(t, deposit.OwnerId, offer.OfferorID)

	CancelOffer(CancelOfferArgs{
		OfferorAddress: address,
		TokenId:        token.TokenID,
		Nonce:          token.Nonce,
		Amount:         "1000000000000000000",
	})
}

func Test_PlaceBid(t *testing.T) {
	connectToDb()
	cache.InitCacher(cfg)

	nonce := uint64(time.Now().Unix())
	token := entities.Token{
		TokenID: "TEST",
		Nonce:   nonce,
	}
	err := storage.AddToken(&token)
	require.Nil(t, err)

	address := "erd12" + fmt.Sprintf("%d", nonce)
	deposit, err := UpdateDeposit(DepositUpdateArgs{
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
	require.Equal(t, entities.ProfferType("Bid"), offer.Type)
	require.Equal(t, token.ID, offer.TokenID)
	require.Equal(t, deposit.OwnerId, offer.OfferorID)
}
