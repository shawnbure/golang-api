package services

import (
	"encoding/json"
	"gorm.io/datatypes"
	"strconv"
	"strings"
	"testing"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
)

func Test_ListAsset(t *testing.T) {
	connectToDb()

	collection := entities.Collection{
		Name:        "col",
		TokenID:     "",
		Description: "",
		CreatorID:   0,
	}
	_ = storage.AddCollection(&collection)

	args := ListAssetArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        13,
		Uri:          "uri",
		Price:        "1000",
		TxHash:       "txHash",
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

	expectedAsset := entities.Asset{
		ID:           asset.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		PriceNominal: 1_000_000_000_000_000_000_000,
		Link:         "uri",
		Listed:       true,
		OwnerId:      ownerAccount.ID,
		CollectionID: asset.CollectionID,
	}
	require.Equal(t, expectedAsset, *asset)
}

func Test_SellAsset(t *testing.T) {
	connectToDb()

	collection := entities.Collection{
		Name:        "col",
		TokenID:     "",
		Description: "",
		CreatorID:   0,
	}
	_ = storage.AddCollection(&collection)

	listArgs := ListAssetArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        13,
		Uri:          "uri",
		Price:        "1000",
		TxHash:       "txHash",
	}
	ListAsset(listArgs)

	buyArgs := BuyAssetArgs{
		OwnerAddress: "ownerAddress",
		BuyerAddress: "buyerAddress",
		TokenId:      "tokenId",
		Nonce:        13,
		Uri:          "col",
		Price:        "1000",
		TxHash:       "txHashBuy",
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

	expectedAsset := entities.Asset{
		ID:           asset.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		PriceNominal: 1_000_000_000_000_000_000_000,
		Link:         "uri",
		Listed:       false,
		OwnerId:      0,
		CollectionID: asset.CollectionID,
	}
	require.Equal(t, expectedAsset, *asset)
}

func Test_WithdrawAsset(t *testing.T) {
	connectToDb()

	collection := entities.Collection{
		Name:        "col",
		TokenID:     "",
		Description: "",
		CreatorID:   0,
	}
	_ = storage.AddCollection(&collection)

	listArgs := ListAssetArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        13,
		Uri:          "uri",
		Price:        "1000",
		TxHash:       "txHash",
	}
	ListAsset(listArgs)

	withdrawArgs := WithdrawAssetArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        13,
		Uri:          "col",
		Price:        "1000",
		TxHash:       "txHashWithdraw",
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

	expectedAsset := entities.Asset{
		ID:           asset.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		PriceNominal: 1_000_000_000_000_000_000_000,
		Link:         "uri",
		Listed:       false,
		OwnerId:      0,
		CollectionID: asset.CollectionID,
	}
	require.Equal(t, expectedAsset, *asset)
}

func Test_GetPriceNominal(T *testing.T) {
	hex := strconv.FormatInt(1_000_000_000_000_000_000, 16)
	priceNominal, err := GetPriceNominal(hex)
	require.Nil(T, err)
	require.Equal(T, priceNominal, float64(1))

	hex = strconv.FormatInt(1_000_000_000_000_000, 16)
	priceNominal, err = GetPriceNominal(hex)
	require.Nil(T, err)
	require.Equal(T, priceNominal, 0.001)

	hex = strconv.FormatInt(100_000_000_000_000, 16)
	priceNominal, err = GetPriceNominal(hex)
	require.Nil(T, err)
	require.Equal(T, priceNominal, float64(0))
}

func Test_GetPriceDenominated(T *testing.T) {
	price := float64(1)
	require.Equal(T, GetPriceDenominated(price).Text(10), "1000000000000000000")

	price = 1000
	require.Equal(T, GetPriceDenominated(price).Text(10), "1000000000000000000000")

	price = 0.001
	require.Equal(T, GetPriceDenominated(price).Text(10), "1000000000000000")

	price = 0.0001
	require.Equal(T, GetPriceDenominated(price).Text(10), "0")
}

func Test_GetAssetLinkResponse(t *testing.T) {
	asset := entities.Asset{
		Nonce: 1,
		Link:  "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t",
	}

	assetLinkWithNonce := GetAssetLinkWithNonce(&asset)
	response, err := HttpGetRaw(assetLinkWithNonce)
	require.Nil(t, err)

	responseLen := len(response)
	require.GreaterOrEqual(t, responseLen, 0)

	attribute := "\"value\":\"Lightning Bolts\",\"trait_type\":\"Earrings\""
	require.True(t, strings.Contains(response, attribute))

	attrs, err := ConstructAttributesJsonFromResponse(response)
	require.Nil(t, err)
	connectToDb()

	attrsMap := make(map[string]string)
	err = json.Unmarshal(*attrs, &attrsMap)
	require.Nil(t, err)
	require.Equal(t, attrsMap["Earrings"], "Lightning Bolts")

	asset.Attributes = *attrs
	err = storage.AddAsset(&asset)
	require.Nil(t, err)

	db, err := storage.GetDBOrError()
	require.Nil(t, err)

	var assetRead entities.Asset
	txRead := db.First(&assetRead, datatypes.JSONQuery("attributes").Equals("Lightning Bolts", "Earrings"))
	require.Nil(t, txRead.Error)
	require.Equal(t, asset.Link, "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t")
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
