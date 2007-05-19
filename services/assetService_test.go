package services

import (
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ListAsset(t *testing.T) {
	connectToDb(t)

	collection := data.Collection{
		Name:        "col",
		TokenID:     "",
		Description: "",
		CreatorID:   0,
	}
	_ = storage.AddNewCollection(&collection)

	ListAsset("ownerAddress", "tokenId", 13, "uri", "col", "1000", "txHash")

	ownerAccount, err := storage.GetAccountByAddress("ownerAddress")
	require.Nil(t, err)
	require.Equal(t, ownerAccount.Address, "ownerAddress")

	transaction, err := storage.GetTransactionByHash("txHash")
	require.Nil(t, err)
	require.Equal(t, transaction.Hash, "txHash")

	asset, err := storage.GetAssetByTokenIdAndNonce("tokenId", 13)
	require.Nil(t, err)

	expectedAsset := data.Asset{
		ID:           asset.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		Price:        "1000",
		Link:         "uri",
		Listed:       true,
		OwnerId:      ownerAccount.ID,
		CollectionID: asset.CollectionID,
	}
	require.Equal(t, expectedAsset, *asset)
}

func Test_SellAsset(t *testing.T) {
	connectToDb(t)

	collection := data.Collection{
		Name:        "col",
		TokenID:     "",
		Description: "",
		CreatorID:   0,
	}
	_ = storage.AddNewCollection(&collection)

	ListAsset("ownerAddress", "tokenId", 13, "uri", "col", "1000", "txHash")
	BuyAsset("ownerAddress", "buyerAddress", "tokenId", 13 , "1000", "txHashBuy")

	ownerAccount, err := storage.GetAccountByAddress("ownerAddress")
	require.Nil(t, err)
	require.Equal(t, ownerAccount.Address, "ownerAddress")

	buyerAccount, err := storage.GetAccountByAddress("buyerAddress")
	require.Nil(t, err)
	require.Equal(t, buyerAccount.Address, "buyerAddress")

	transaction, err := storage.GetTransactionByHash("txHashBuy")
	require.Nil(t, err)
	require.Equal(t, transaction.Hash, "txHashBuy")

	asset, err := storage.GetAssetByTokenIdAndNonce("tokenId", 13)
	require.Nil(t, err)

	expectedAsset := data.Asset{
		ID:           asset.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		Price:        "1000",
		Link:         "uri",
		Listed:       false,
		OwnerId:      0,
		CollectionID: asset.CollectionID,
	}
	require.Equal(t, expectedAsset, *asset)
}

func Test_WithdrawAsset(t *testing.T) {
	connectToDb(t)

	collection := data.Collection{
		Name:        "col",
		TokenID:     "",
		Description: "",
		CreatorID:   0,
	}
	_ = storage.AddNewCollection(&collection)

	ListAsset("ownerAddress", "tokenId", 13, "uri", "col", "1000", "txHash")
	WithdrawAsset("ownerAddress", "tokenId", 13 , "1000", "txHashWithdraw")

	ownerAccount, err := storage.GetAccountByAddress("ownerAddress")
	require.Nil(t, err)
	require.Equal(t, ownerAccount.Address, "ownerAddress")

	transaction, err := storage.GetTransactionByHash("txHashWithdraw")
	require.Nil(t, err)
	require.Equal(t, transaction.Hash, "txHashWithdraw")

	asset, err := storage.GetAssetByTokenIdAndNonce("tokenId", 13)
	require.Nil(t, err)

	expectedAsset := data.Asset{
		ID:           asset.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		Price:        "1000",
		Link:         "uri",
		Listed:       false,
		OwnerId:      0,
		CollectionID: asset.CollectionID,
	}
	require.Equal(t, expectedAsset, *asset)
}
