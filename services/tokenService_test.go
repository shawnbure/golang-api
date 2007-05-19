package services

import (
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
)

func Test_ListToken(t *testing.T) {
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

	token, err := storage.GetTokenByTokenIdAndNonce("tokenId", 13)
	require.Nil(t, err)

	expectedToken := entities.Token{
		ID:           token.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		PriceNominal: 1_000,
		Status:       entities.List,
		OwnerId:      ownerAccount.ID,
		CollectionID: token.CollectionID,
	}
	require.Equal(t, expectedToken, *token)
}

func Test_SellToken(t *testing.T) {
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

	token, err := storage.GetTokenByTokenIdAndNonce("tokenId", 13)
	require.Nil(t, err)

	expectedToken := entities.Token{
		ID:           token.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		PriceNominal: 1_000,
		Status:       entities.List,
		OwnerId:      0,
		CollectionID: token.CollectionID,
	}
	require.Equal(t, expectedToken, *token)
}

func Test_WithdrawToken(t *testing.T) {
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

	token, err := storage.GetTokenByTokenIdAndNonce("tokenId", 13)
	require.Nil(t, err)

	expectedToken := entities.Token{
		ID:           token.ID,
		TokenID:      "tokenId",
		Nonce:        13,
		PriceNominal: 1_000,
		Status:       entities.List,
		OwnerId:      0,
		CollectionID: token.CollectionID,
	}
	require.Equal(t, expectedToken, *token)
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

func Test_GetTokenLinkResponse(t *testing.T) {
	token := entities.Token{
		Nonce:        1,
		MetadataLink: "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t",
	}

	osResponse, err := GetOSMetadataForToken(token.MetadataLink, token.Nonce)
	require.Nil(t, err)

	attrs, err := ConstructAttributesJsonFromResponse(osResponse)
	require.Nil(t, err)
	connectToDb()

	attrsMap := make(map[string]string)
	err = json.Unmarshal(*attrs, &attrsMap)
	require.Nil(t, err)
	require.Equal(t, attrsMap["Earrings"], "Lightning Bolts")

	token.Attributes = *attrs
	err = storage.AddToken(&token)
	require.Nil(t, err)

	db, err := storage.GetDBOrError()
	require.Nil(t, err)

	var tokenRead entities.Token
	txRead := db.First(&tokenRead, datatypes.JSONQuery("attributes").Equals("Lightning Bolts", "Earrings"))
	require.Nil(t, txRead.Error)
	require.Equal(t, token.MetadataLink, "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t")
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

	token, err := storage.GetTokenByTokenIdAndNonce("tokenId", nonce)
	require.Nil(t, err)
	require.True(t, strings.Contains(token.ImageLink, "ipfs://"))
	require.Equal(t, token.MetadataLink, "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t")
	require.True(t, strings.Contains(string(token.Attributes), "Lips"))
	require.True(t, strings.Contains(token.TokenName, "Woman"))

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

	token, err = storage.GetTokenByTokenIdAndNonce("tokenId", nonce)
	require.Nil(t, err)
	require.Equal(t, token.ImageLink, "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t")
	require.Equal(t, token.MetadataLink, "")
	require.Equal(t, string(token.Attributes), `{"ceva": "ceva"}`)
	require.Equal(t, token.TokenName, "some_name")
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

func Test_StartAuction(t *testing.T) {
	connectToDb()
	cache.InitCacher(cfg)

	nonce := uint64(time.Now().Unix())
	address := "erd12" + fmt.Sprintf("%d", nonce)
	token, err := StartAuction(StartAuctionArgs{
		OwnerAddress:     address,
		Nonce:            nonce,
		FirstLink:        "abcdef",
		MinBid:           "1000000000000000000",
		StartTime:        10,
		Deadline:         110,
		RoyaltiesPercent: 10,
	})
	require.Nil(t, err)
	require.NotEqual(t, uint64(0), token.ID)

	owner, err := storage.GetAccountByAddress(address)
	require.Nil(t, err)
	require.NotEqual(t, uint64(0), owner.ID)

	require.Equal(t, nonce, token.Nonce)
	require.Equal(t, "abcdef", token.ImageLink)
	require.Equal(t, "1000000000000000000", token.PriceString)
	require.Equal(t, entities.TokenStatus(entities.Auction), token.Status)
	require.Equal(t, owner.ID, token.OwnerId)
}

func Test_StartAuctionEndAuction(t *testing.T) {
	connectToDb()
	cache.InitCacher(cfg)

	nonce := uint64(time.Now().Unix())
	address := "erd12" + fmt.Sprintf("%d", nonce)
	token, err := StartAuction(StartAuctionArgs{
		OwnerAddress:     address,
		Nonce:            nonce,
		FirstLink:        "abcdef",
		MinBid:           "1000000000000000000",
		StartTime:        10,
		Deadline:         110,
		RoyaltiesPercent: 10,
	})
	require.Nil(t, err)
	require.NotEqual(t, uint64(0), token.ID)

	owner, err := storage.GetAccountByAddress(address)
	require.Nil(t, err)
	require.NotEqual(t, uint64(0), owner.ID)

	require.Equal(t, nonce, token.Nonce)
	require.Equal(t, "abcdef", token.ImageLink)
	require.Equal(t, "1000000000000000000", token.PriceString)
	require.Equal(t, entities.TokenStatus(entities.Auction), token.Status)
	require.Equal(t, owner.ID, token.OwnerId)

	EndAuction(EndAuctionArgs{
		TokenId:   token.TokenID,
		Nonce:     nonce,
		Winner:    address,
		Amount:    "1000000000000000000",
	})

	tokenAfterEnd, err := storage.GetTokenByTokenIdAndNonce(token.TokenID, token.Nonce)
	require.Nil(t, err)
	require.Equal(t, uint64(0), tokenAfterEnd.OwnerId)
	require.Equal(t, 4722.366, tokenAfterEnd.LastBuyPriceNominal)
	require.Equal(t, entities.TokenStatus(entities.None), tokenAfterEnd.Status)
}
