package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/datatypes"
	"math/big"
	"strconv"
	"strings"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/boltdb/bolt"
	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/stats/collstats"
	"github.com/erdsea/erdsea-api/storage"
)

type AvailableTokensRequest struct {
	Tokens []string `json:"tokens"`
}

type AvailableToken struct {
	Collection struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"collection"`
	Token struct {
		Id        string `json:"id"`
		Nonce     uint64 `json:"nonce"`
		Name      string `json:"name"`
		Available bool   `json:"available"`
	} `json:"token"`
}

type AvailableTokensResponse struct {
	Tokens map[string]AvailableToken `json:"tokens"`
}

type TokenLinkResponse struct {
	Name       string           `json:"name"`
	Image      string           `json:"image"`
	Attributes []dtos.Attribute `json:"attributes"`
}

type TokenCacheInfo struct {
	TokenDbId uint64
	TokenName string
}

const (
	minPriceUnit            = 1000
	minPercentUnit          = 1000
	minPercentRoyaltiesUnit = 100
	minPriceDecimals        = 15

	maxTokenLinkResponseSize = 1024
)

var (
	TokenIdToDbIdCacheInfo = []byte("tokenToId")

	baseExp = big.NewInt(10)

	log = logger.GetOrCreate("services")
)

func ListToken(args ListTokenArgs) {
	priceNominal, err := GetPriceNominal(args.Price)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return
	}

	ownerAccount, err := GetOrCreateAccount(args.OwnerAddress)
	if err != nil {
		log.Debug("could not get or create account", "err", err)
		return
	}

	collectionId := uint64(0)
	collection, err := storage.GetCollectionByTokenId(args.TokenId)
	if err == nil {
		collectionId = collection.ID
	}

	token, err := storage.GetTokenByTokenIdAndNonce(args.TokenId, args.Nonce)

	var innerErr error
	if err != nil {
		newToken := ConstructNewTokenFromListArgs(args)
		token = &newToken
		token.Listed = true
		token.PriceString = args.Price
		token.PriceNominal = priceNominal
		token.OwnerId = ownerAccount.ID
		token.CollectionID = collectionId
		innerErr = storage.AddToken(token)

		_, cacheErr := AddTokenToCache(token.TokenID, token.Nonce, token.TokenName, token.ID)
		if cacheErr != nil {
			log.Error("could not add token to cache")
		}
	} else {
		token.Listed = true
		token.PriceString = args.Price
		token.PriceNominal = priceNominal
		token.OwnerId = ownerAccount.ID
		token.CollectionID = collectionId
		innerErr = storage.UpdateToken(token)
	}

	if innerErr != nil {
		log.Debug("could not create or update token", "err", innerErr)
		return
	}

	transaction := entities.Transaction{
		Hash:         args.TxHash,
		Type:         entities.ListToken,
		PriceNominal: priceNominal,
		Timestamp:    args.Timestamp,
		SellerID:     ownerAccount.ID,
		BuyerID:      0,
		TokenID:      token.ID,
		CollectionID: collectionId,
	}

	AddTransaction(&transaction)
}

func BuyToken(args BuyTokenArgs) {
	priceNominal, err := GetPriceNominal(args.Price)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return
	}

	ownerAccount, err := storage.GetAccountByAddress(args.OwnerAddress)
	if err != nil {
		log.Debug("could not get owner account", "err", err)
		return
	}

	buyerAccount, err := GetOrCreateAccount(args.BuyerAddress)
	if err != nil {
		log.Debug("could not get or create account", "err", err)
		return
	}

	token, err := storage.GetTokenByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("could not get token", "err", err)
		return
	}

	token.Listed = false
	// This was to be reset since the token will no longer be on the marketplace.
	// Could have been kept like this, but bugs may appear when querying.
	token.OwnerId = 0
	err = storage.UpdateToken(token)
	if err != nil {
		log.Debug("could not update token", "err", err)
		return
	}

	transaction := entities.Transaction{
		Hash:         args.TxHash,
		Type:         entities.BuyToken,
		PriceNominal: priceNominal,
		Timestamp:    args.Timestamp,
		SellerID:     ownerAccount.ID,
		BuyerID:      buyerAccount.ID,
		TokenID:      token.ID,
		CollectionID: token.CollectionID,
	}

	AddTransaction(&transaction)
}

func WithdrawToken(args WithdrawTokenArgs) {
	priceNominal, err := GetPriceNominal(args.Price)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return
	}

	ownerAccount, err := storage.GetAccountByAddress(args.OwnerAddress)
	if err != nil {
		log.Debug("could not get owner account", err)
		return
	}

	token, err := storage.GetTokenByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("could not get token", "err", err)
		return
	}

	token.Listed = false
	// This was to be reset since the token will no longer be on the marketplace.
	// Could have been kept like this, but bugs may appear when trying when querying.
	token.OwnerId = 0
	err = storage.UpdateToken(token)
	if err != nil {
		log.Debug("could not update token", "err", err)
		return
	}

	transaction := entities.Transaction{
		Hash:         args.TxHash,
		Type:         entities.WithdrawToken,
		PriceNominal: priceNominal,
		Timestamp:    args.Timestamp,
		SellerID:     0,
		BuyerID:      ownerAccount.ID,
		TokenID:      token.ID,
		CollectionID: token.CollectionID,
	}

	AddTransaction(&transaction)
}

func GetExtendedTokenData(tokenId string, nonce uint64) (*dtos.ExtendedTokenDto, error) {
	token, err := storage.GetTokenByTokenIdAndNonce(tokenId, nonce)
	if err != nil {
		return nil, err
	}

	collection, err := storage.GetCollectionById(token.CollectionID)
	if err != nil {
		return nil, err
	}

	owner, err := storage.GetAccountById(token.OwnerId)
	if err != nil {
		return nil, err
	}

	creator, err := storage.GetAccountById(collection.CreatorID)
	if err != nil {
		return nil, err
	}

	collStats, err := collstats.GetStatisticsForTokenId(tokenId)
	if err != nil {
		collStats = &dtos.CollectionStatistics{}
	}

	return dtos.CreateExtendedTokenDto(
		*token,
		*collection,
		owner.Name,
		owner.Address,
		creator.Name,
		creator.Address,
		*collStats,
	)
}

func ConstructNewTokenFromListArgs(args ListTokenArgs) entities.Token {
	token := entities.Token{
		TokenID:          args.TokenId,
		Nonce:            args.Nonce,
		RoyaltiesPercent: GetRoyaltiesPercentNominal(args.RoyaltiesPercent),
		MetadataLink:     "",
		CreatedAt:        args.Timestamp,
		Listed:           true,
		Attributes:       datatypes.JSON(""),
		TokenName:        "",
		ImageLink:        "",
		Hash:             args.Hash,
	}

	osResponse, err := GetOSMetadataForToken(args.LastLink, args.Nonce)
	if err == nil {
		token.MetadataLink = args.LastLink
		token.TokenName = osResponse.Name
		token.ImageLink = osResponse.Image

		attributesJson, innerErr := ConstructAttributesJsonFromResponse(osResponse)
		if innerErr != nil {
			log.Debug("could not parse os response for attributes", "link", args.LastLink)
		} else {
			token.Attributes = *attributesJson
		}
	} else {
		token.TokenName = args.TokenName
		token.ImageLink = args.FirstLink
	}

	return token
}

func GetOSMetadataForToken(link string, nonce uint64) (*TokenLinkResponse, error) {
	var response TokenLinkResponse

	link = link + "/" + strconv.FormatUint(nonce, 10)
	responseRaw, err := HttpGetRaw(link)
	if err != nil {
		return nil, err
	}
	if len(responseRaw) > maxTokenLinkResponseSize {
		return nil, errors.New("response too long")
	}

	err = json.Unmarshal([]byte(responseRaw), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func ConstructAttributesJsonFromResponse(response *TokenLinkResponse) (*datatypes.JSON, error) {
	attrsBytes, err := json.Marshal(response.Attributes)
	if err != nil {
		return nil, err
	}

	attrsJson := datatypes.JSON(attrsBytes)
	return &attrsJson, err
}

func GetPriceNominal(priceHex string) (float64, error) {
	priceBigUint, success := big.NewInt(0).SetString(priceHex, 16)
	if !success {
		return 0, errors.New("could not parse price")
	}

	denominatorBigUint := big.NewInt(0).Exp(baseExp, big.NewInt(minPriceDecimals), nil)
	priceNominalInt := big.NewInt(0).Div(priceBigUint, denominatorBigUint).Int64()
	priceNominal := float64(priceNominalInt) / minPercentUnit
	return priceNominal, nil
}

func GetPriceDenominated(price float64) *big.Int {
	priceInt := int64(price * minPriceUnit)
	if priceInt <= 0 {
		log.Error("price less than min threshold",
			"min_threshold_multiplied", "1",
			"min_threshold_nominal", 1/minPriceUnit,
			"price_int", priceInt,
		)
	}

	denominatorBigUint := big.NewInt(0).Exp(baseExp, big.NewInt(minPriceDecimals), nil)

	priceBigUint := big.NewInt(0).Mul(big.NewInt(priceInt), denominatorBigUint)
	return priceBigUint
}

func GetRoyaltiesPercentNominal(percent uint64) float64 {
	return float64(percent) / minPercentRoyaltiesUnit
}

func AddTransaction(tx *entities.Transaction) {
	err := storage.AddTransaction(tx)
	if err != nil {
		log.Debug("could not create new transaction", "err", err)
		return
	}
}

func AddTokenToCache(tokenId string, nonce uint64, tokenName string, tokenDbId uint64) (*TokenCacheInfo, error) {
	db := cache.GetBolt()
	cacheInfo := TokenCacheInfo{
		TokenDbId: tokenDbId,
		TokenName: tokenName,
	}

	entryBytes, err := json.Marshal(&cacheInfo)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, innerErr := tx.CreateBucketIfNotExists(TokenIdToDbIdCacheInfo)
		if innerErr != nil {
			return innerErr
		}

		key := fmt.Sprintf("%s-%d", tokenId, nonce)
		innerErr = bucket.Put([]byte(key), entryBytes)
		return innerErr
	})

	return &cacheInfo, nil
}

func GetTokenCacheInfo(tokenId string, nonce uint64) (*TokenCacheInfo, error) {
	db := cache.GetBolt()

	var bytes []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(TokenIdToDbIdCacheInfo)
		if bucket == nil {
			return errors.New("no bucket for token cache")
		}

		key := fmt.Sprintf("%s-%d", tokenId, nonce)
		bytes = bucket.Get([]byte(key))
		return nil
	})
	if err != nil {
		return nil, err
	}

	var cacheInfo TokenCacheInfo
	err = json.Unmarshal(bytes, &cacheInfo)
	if err != nil {
		return nil, err
	}

	return &cacheInfo, nil
}

func GetOrAddTokenCacheInfo(tokenId string, nonce uint64) (*TokenCacheInfo, error) {
	cacheInfo, err := GetTokenCacheInfo(tokenId, nonce)
	if err != nil {
		token, innerErr := storage.GetTokenByTokenIdAndNonce(tokenId, nonce)
		if innerErr != nil {
			return nil, innerErr
		}

		cacheInfo, innerErr = AddTokenToCache(tokenId, nonce, token.TokenName, token.ID)
		if innerErr != nil {
			return nil, innerErr
		}
	}

	return cacheInfo, nil
}

func GetAvailableTokens(args AvailableTokensRequest) AvailableTokensResponse {
	var response AvailableTokensResponse
	response.Tokens = make(map[string]AvailableToken)

	for _, token := range args.Tokens {
		parts := strings.Split(token, "-")
		if len(parts) != 3 {
			continue
		}

		tokenId := parts[0] + "-" + parts[1]
		nonce, err := strconv.ParseUint(parts[2], 16, 64)
		if err != nil {
			continue
		}

		tokenName := ""
		tokenAvailable := false
		tokenCacheInfo, err := GetOrAddTokenCacheInfo(tokenId, nonce)
		if err == nil {
			tokenAvailable = true
			tokenName = tokenCacheInfo.TokenName
		}

		collectionName := ""
		collectionCacheInfo, err := collstats.GetOrAddCollectionCacheInfo(tokenId)
		if err == nil {
			collectionName = collectionCacheInfo.CollectionName
		}

		response.Tokens[token] = AvailableToken{
			Collection: struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			}{
				Id:   tokenId,
				Name: collectionName,
			},
			Token: struct {
				Id        string `json:"id"`
				Nonce     uint64 `json:"nonce"`
				Name      string `json:"name"`
				Available bool   `json:"available"`
			}{
				Id:        tokenId,
				Nonce:     nonce,
				Name:      tokenName,
				Available: tokenAvailable,
			},
		}
	}

	return response
}
