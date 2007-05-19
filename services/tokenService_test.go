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
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
)

var blockchainCfg = config.BlockchainConfig{
	ProxyUrl:           "https://devnet-gateway.elrond.com",
	MarketplaceAddress: "erd1qqqqqqqqqqqqqpgqm4dmwyxc5fsj49z3jcu9h08azjrcf60kt9uspxs483",
}

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
	ListToken(args, blockchainCfg.ProxyUrl, blockchainCfg.MarketplaceAddress)

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
	ListToken(listArgs, blockchainCfg.ProxyUrl, blockchainCfg.MarketplaceAddress)

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
	ListToken(listArgs, blockchainCfg.ProxyUrl, blockchainCfg.MarketplaceAddress)

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
		MetadataLink: "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t/1",
	}

	attrs := GetAttributesFromMetadata(token.MetadataLink)
	require.NotZero(t, len(attrs))

	attrsMap := make(map[string]string)
	err := json.Unmarshal(attrs, &attrsMap)
	require.Nil(t, err)
	require.Equal(t, attrsMap["Earrings"], "Lightning Bolts")

	connectToDb()
	token.Attributes = attrs
	err = storage.AddToken(&token)
	require.Nil(t, err)

	db, err := storage.GetDBOrError()
	require.Nil(t, err)

	var tokenRead entities.Token
	txRead := db.First(&tokenRead, datatypes.JSONQuery("attributes").Equals("Lightning Bolts", "Earrings"))
	require.Nil(t, txRead.Error)
	require.Equal(t, "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t/1", token.MetadataLink)
}

func Test_ErdCompatibility(t *testing.T) {
	connectToDb()
	cache.InitCacher(cacheCfg)

	nonce := uint64(69)
	listArgs := ListTokenArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        nonce,
		TokenName:    "some_name",
		Price:        "3635C9ADC5DEA00000",
		Attributes:   "some_attrs",
		FirstLink:    "some_link",
		SecondLink:   "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t/69",
	}
	ListToken(listArgs, blockchainCfg.ProxyUrl, blockchainCfg.MarketplaceAddress)

	token, err := storage.GetTokenByTokenIdAndNonce("tokenId", nonce)
	require.Nil(t, err)
	require.Equal(t, "some_link", token.ImageLink)
	require.Equal(t, token.MetadataLink, "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t/69")
	require.True(t, strings.Contains(string(token.Attributes), "Lips"))
	require.Equal(t, "some_name", token.TokenName)

	nonce = nonce + 1
	listArgs = ListTokenArgs{
		OwnerAddress: "ownerAddress",
		TokenId:      "tokenId",
		Nonce:        nonce,
		TokenName:    "some_name",
		Price:        "3635C9ADC5DEA00000",
		Attributes:   `{"ceva": "ceva"}`,
		FirstLink:    "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t/70",
		SecondLink:   "some_link",
	}
	ListToken(listArgs, blockchainCfg.ProxyUrl, blockchainCfg.MarketplaceAddress)

	token, err = storage.GetTokenByTokenIdAndNonce("tokenId", nonce)
	require.Nil(t, err)
	require.Equal(t, token.ImageLink, "https://wow-prod-nftribe.s3.eu-west-2.amazonaws.com/t/70")
	require.Equal(t, token.MetadataLink, "some_link")
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
	cache.InitCacher(cacheCfg)

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
	}, blockchainCfg.ProxyUrl, blockchainCfg.MarketplaceAddress)
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
	cache.InitCacher(cacheCfg)

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
	}, blockchainCfg.ProxyUrl, blockchainCfg.MarketplaceAddress)
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
		TokenId: token.TokenID,
		Nonce:   nonce,
		Winner:  address,
		Amount:  "1000000000000000000",
	})

	tokenAfterEnd, err := storage.GetTokenByTokenIdAndNonce(token.TokenID, token.Nonce)
	require.Nil(t, err)
	require.Equal(t, uint64(0), tokenAfterEnd.OwnerId)
	require.Equal(t, 4722.366, tokenAfterEnd.LastBuyPriceNominal)
	require.Equal(t, entities.TokenStatus(entities.None), tokenAfterEnd.Status)
}

func Test_GetMetadata(t *testing.T) {
	cache.InitCacher(cacheCfg)

	link := "https://www.chubbiverse.com/api/meta/1/1"
	bytes, err := TryGetResponseCached(link)
	require.Nil(t, err)
	println(bytes)
	println(len(bytes))
	var response dtos.MetadataLinkResponse
	err = json.Unmarshal([]byte(bytes), &response)
	require.Nil(t, err)

	println("--------------------")

	link = "https://gateway.pinata.cloud/ipfs/QmPxQivTP7tncEkyrB7yLG7XmQXXsj8H9GfZ7vdH69eJ8k/1"
	bytes, err = TryGetResponseCached(link)
	require.Nil(t, err)
	println(bytes)
	println(len(bytes))
	err = json.Unmarshal([]byte(bytes), &response)
	require.Nil(t, err)

	println("--------------------")

	link = "https://gateway.pinata.cloud/ipfs/QmZRb9AxuDAe8KfcVGw6XUQa1bH1pRPT2hbTWwWb48tLgb/1"
	bytes, err = TryGetResponseCached(link)
	require.Nil(t, err)
	println(bytes)
	println(len(bytes))
	err = json.Unmarshal([]byte(bytes), &response)
	require.Nil(t, err)

	println("--------------------")

	link = "https://nftartisans.mypinata.cloud/ipfs/QmRR919iHG8frctbLhktKzSgWJyCGUp8rktJaczDzyvAke/1"
	bytes, err = TryGetResponseCached(link)
	require.Nil(t, err)
	println(bytes)
	println(len(bytes))
	err = json.Unmarshal([]byte(bytes), &response)
	require.Nil(t, err)

	println("--------------------")

	link = "https://galacticapes.mypinata.cloud/ipfs/QmRDdnGJYQhPjq8ebxtsea8cptCUkgZwgBMnq1grW8mWJr/1"
	bytes, err = TryGetResponseCached(link)
	require.Nil(t, err)
	println(bytes)
	println(len(bytes))
	err = json.Unmarshal([]byte(bytes), &response)
	require.Nil(t, err)


	println("--------------------")
}

func Test_RefreshMetadata(t *testing.T) {
	connectToDb()
	cache.InitCacher(cacheCfg)

	token := entities.Token{
		TokenID:             "CHUB-dc4906",
		Nonce:               5,
		MetadataLink:        "",
	}

	err := storage.AddToken(&token)
	require.Nil(t, err)

	attrs, err := RefreshMetadata(blockchainCfg.ProxyUrl, &token, "erd17s2pz8qrds6ake3qwheezgy48wzf7dr5nhdpuu2h4rr4mt5rt9ussj7xzh")
	require.Nil(t, err)
	require.NotEqual(t, 0, len(attrs))

	readToken, err := storage.GetTokenById(token.ID)
	require.Nil(t, err)

	var tokenAttrs map[string]interface{}
	err = json.Unmarshal(token.Attributes, &tokenAttrs)
	require.Nil(t, err)


	var readTokenAttrs map[string]interface{}
	err = json.Unmarshal(readToken.Attributes, &readTokenAttrs)
	require.Nil(t, err)

	require.Equal(t, tokenAttrs, readTokenAttrs)
}
