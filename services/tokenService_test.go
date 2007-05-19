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

	args := ListTokenArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        13,
		Price:        "3635C9ADC5DEA00000",
		TxHash:       "txHash",
	}
	ListToken(args)

	ownerAccount, err := storage.GetAccountByAddress("ownerAddress")
	require.Nil(t, err)
	require.Equal(t, ownerAccount.Address, "ownerAddress")

	transaction, err := storage.GetTransactionByHash("txHash")
	require.Nil(t, err)
	require.Equal(t, transaction.Hash, "txHash")

	asset, err := storage.GetTokenByTokenIdAndNonce("tokenId", 13)
	require.Nil(t, err)

	expectedAsset := entities.Token{
		ID:           asset.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		PriceNominal: 1_000,
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

	listArgs := ListTokenArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        13,
		Price:        "3635C9ADC5DEA00000",
		TxHash:       "txHash",
	}
	ListToken(listArgs)

	buyArgs := BuyTokenArgs{
		OwnerAddress: "ownerAddress",
		BuyerAddress: "buyerAddress",
		TokenId:      "tokenId",
		Nonce:        13,
		Price:        "3635C9ADC5DEA00000",
		TxHash:       "txHashBuy",
	}
	BuyToken(buyArgs)

	ownerAccount, err := storage.GetAccountByAddress("ownerAddress")
	require.Nil(t, err)
	require.Equal(t, ownerAccount.Address, "ownerAddress")

	buyerAccount, err := storage.GetAccountByAddress("buyerAddress")
	require.Nil(t, err)
	require.Equal(t, buyerAccount.Address, "buyerAddress")

	transaction, err := storage.GetTransactionByHash("txHashBuy")
	require.Nil(t, err)
	require.Equal(t, transaction.Hash, "txHashBuy")

	asset, err := storage.GetTokenByTokenIdAndNonce("tokenId", 13)
	require.Nil(t, err)

	expectedAsset := entities.Token{
		ID:           asset.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		PriceNominal: 1_000,
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

	listArgs := ListTokenArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        13,
		Price:        "3635C9ADC5DEA00000",
		TxHash:       "txHash",
	}
	ListToken(listArgs)

	withdrawArgs := WithdrawTokenArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        13,
		Price:        "3635C9ADC5DEA00000",
		TxHash:       "txHashWithdraw",
	}
	WithdrawToken(withdrawArgs)

	ownerAccount, err := storage.GetAccountByAddress("ownerAddress")
	require.Nil(t, err)
	require.Equal(t, ownerAccount.Address, "ownerAddress")

	transaction, err := storage.GetTransactionByHash("txHashWithdraw")
	require.Nil(t, err)
	require.Equal(t, transaction.Hash, "txHashWithdraw")

	asset, err := storage.GetTokenByTokenIdAndNonce("tokenId", 13)
	require.Nil(t, err)

	expectedAsset := entities.Token{
		ID:           asset.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		PriceNominal: 1_000,
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
	asset := entities.Token{
		Nonce:        1,
		MetadataLink: "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t",
	}

	osResponse, err := GetOSMetadataForToken(asset.MetadataLink, asset.Nonce)
	require.Nil(t, err)

	attrs, err := ConstructAttributesJsonFromResponse(osResponse)
	require.Nil(t, err)
	connectToDb()

	attrsMap := make(map[string]string)
	err = json.Unmarshal(*attrs, &attrsMap)
	require.Nil(t, err)
	require.Equal(t, attrsMap["Earrings"], "Lightning Bolts")

	asset.Attributes = *attrs
	err = storage.AddToken(&asset)
	require.Nil(t, err)

	db, err := storage.GetDBOrError()
	require.Nil(t, err)

	var assetRead entities.Token
	txRead := db.First(&assetRead, datatypes.JSONQuery("attributes").Equals("Lightning Bolts", "Earrings"))
	require.Nil(t, txRead.Error)
	require.Equal(t, asset.MetadataLink, "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t")
}

func Test_ErdCompatibility(t *testing.T) {
	connectToDb()

	nonce := uint64(69)
	listArgs := ListTokenArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        nonce,
		TokenName:    "some_name",
		Price:        "3635C9ADC5DEA00000",
		Attributes:   "some_attrs",
		FirstLink:    "some_link",
		LastLink:     "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t",
	}
	ListToken(listArgs)

	asset, err := storage.GetTokenByTokenIdAndNonce("tokenId", nonce)
	require.Nil(t, err)
	require.True(t, strings.Contains(asset.ImageLink, "ipfs://"))
	require.Equal(t, asset.MetadataLink, "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t")
	require.True(t, strings.Contains(string(asset.Attributes), "Lips"))
	require.True(t, strings.Contains(asset.TokenName, "Woman"))

	nonce = nonce + 1
	listArgs = ListTokenArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        nonce,
		TokenName:    "some_name",
		Price:        "3635C9ADC5DEA00000",
		Attributes:   `{"ceva": "ceva"}`,
		FirstLink:    "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t",
		LastLink:     "some_link",
	}
	ListToken(listArgs)

	asset, err = storage.GetTokenByTokenIdAndNonce("tokenId", nonce)
	require.Nil(t, err)
	require.Equal(t, asset.ImageLink, "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t")
	require.Equal(t, asset.MetadataLink, "")
	require.Equal(t, string(asset.Attributes), `{"ceva": "ceva"}`)
	require.Equal(t, asset.TokenName, "some_name")
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
