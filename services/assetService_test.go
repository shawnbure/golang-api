package services

import (
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ListAsset(t *testing.T) {
	connectToDb()

	collection := data.Collection{
		Name:        "col",
		TokenID:     "",
		Description: "",
		CreatorID:   0,
	}
	_ = storage.AddNewCollection(&collection)

	args := ListAssetArgs{
		OwnerAddress:   "ownerAddress",
		TokenId:        "tokenId",
		Nonce:          13,
		Uri:            "uri",
		Price:          "1000",
		TxHash:         "txHash",
	}
	ListAsset(args)

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
	connectToDb()

	collection := data.Collection{
		Name:        "col",
		TokenID:     "",
		Description: "",
		CreatorID:   0,
	}
	_ = storage.AddNewCollection(&collection)

	listArgs := ListAssetArgs{
		OwnerAddress:   "ownerAddress",
		TokenId:        "tokenId",
		Nonce:          13,
		Uri:            "uri",
		Price:          "1000",
		TxHash:         "txHash",
	}
	ListAsset(listArgs)

	buyArgs := BuyAssetArgs{
		OwnerAddress:   "ownerAddress",
		BuyerAddress:   "buyerAddress",
		TokenId:        "tokenId",
		Nonce:          13,
		Uri:            "col",
		Price:          "1000",
		TxHash:         "txHashBuy",
	}
	BuyAsset(buyArgs)

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
	connectToDb()

	collection := data.Collection{
		Name:        "col",
		TokenID:     "",
		Description: "",
		CreatorID:   0,
	}
	_ = storage.AddNewCollection(&collection)

	listArgs := ListAssetArgs{
		OwnerAddress:   "ownerAddress",
		TokenId:        "tokenId",
		Nonce:          13,
		Uri:            "uri",
		Price:          "1000",
		TxHash:         "txHash",
	}
	ListAsset(listArgs)

	withdrawArgs := WithdrawAssetArgs{
		OwnerAddress:   "ownerAddress",
		TokenId:        "tokenId",
		Nonce:          13,
		Uri:            "col",
		Price:          "1000",
		TxHash:         "txHashWithdraw",
	}
	WithdrawAsset(withdrawArgs)

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

func connectToDb() {
	storage.Connect(config.DatabaseConfig{
		Dialect:       "postgres",
		Host:          "localhost",
		Port:          5432,
		DbName:        "erdsea_db_test",
		User:          "postgres",
		Password:      "root",
		SslMode:       "disable",
		MaxOpenConns:  50,
		MaxIdleConns:  10,
		ShouldMigrate: true,
	})
}
