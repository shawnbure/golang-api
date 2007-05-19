package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
)

const (
	TokenIdMaxLen            = 15
	MaxNameLen               = 20
	MaxLinkLen               = 100
	MaxDescLen               = 1000
	RegisteredNFTsBaseFormat = "%s/address/%s/registered-nfts"
	HttpResponseExpirePeriod = 10 * time.Minute
	StatisticsCacheKeyFormat = "StatisticsForId%d"
	StatisticsExpirePeriod   = 15 * time.Minute
)

type CreateCollectionRequest struct {
	UserAddress   string `json:"userAddress"`
	Name          string `json:"collectionName"`
	TokenId       string `json:"tokenId"`
	Description   string `json:"description"`
	Website       string `json:"website"`
	DiscordLink   string `json:"discordLink"`
	TwitterLink   string `json:"twitterLink"`
	InstagramLink string `json:"instagramLink"`
	TelegramLink  string `json:"telegramLink"`
}

type CollectionStatistics struct {
	ItemsCount   uint64  `json:"itemsCount"`
	OwnersCount  uint64  `json:"ownersCount"`
	FloorPrice   float64 `json:"floorPrice"`
	VolumeTraded float64 `json:"volumeTraded"`
}

type ProxyRegisteredNFTsResponse struct {
	Data struct {
		Tokens []string `json:"tokens"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

func CreateCollection(request *CreateCollectionRequest, blockchainProxy string) error {
	err := checkValidInput(request)
	if err != nil {
		return err
	}

	tokensRegisteredByUser, err := getTokensRegisteredByUser(request.UserAddress, blockchainProxy)
	if err != nil {
		return err
	}
	if !contains(tokensRegisteredByUser, request.TokenId) {
		return errors.New("token not owner by user")
	}

	_, err = storage.GetCollectionByName(request.Name)
	if err == nil {
		return errors.New("collection name already taken")
	}

	_, err = storage.GetCollectionByTokenId(request.TokenId)
	if err == nil {
		return errors.New("token id already has a collection associated")
	}

	account, err := GetOrCreateAccount(request.UserAddress)
	if err != nil {
		return err
	}

	collection := &data.Collection{
		ID:            0,
		Name:          request.Name,
		TokenID:       request.TokenId,
		Description:   request.Description,
		Website:       request.Website,
		DiscordLink:   request.DiscordLink,
		TwitterLink:   request.TwitterLink,
		InstagramLink: request.InstagramLink,
		TelegramLink:  request.TelegramLink,
		CreatorID:     account.ID,
		CreatedAt:     uint64(time.Now().Unix()),
	}

	return storage.AddNewCollection(collection)
}

func GetStatisticsForCollection(collectionId uint64) (*CollectionStatistics, error) {
	var stats CollectionStatistics
	cacheKey := fmt.Sprintf(StatisticsCacheKeyFormat, collectionId)

	err := cache.GetCacher().Get(cacheKey, &stats)
	if err == nil {
		return &stats, nil
	}

	numItems, err := storage.CountListedAssetsByCollectionId(collectionId)
	if err != nil {
		return nil, err
	}

	numOwners, err := storage.CountUniqueOwnersWithListedAssetsByCollectionId(collectionId)
	if err != nil {
		return nil, err
	}

	//TODO: refactor this to something smarter. Min price is not good
	minPrice, err := storage.GetMinBuyPriceForTransactionsWithCollectionId(collectionId)
	if err != nil {
		return nil, err
	}

	sumPrice, err := storage.GetSumBuyPriceForTransactionsWithCollectionId(collectionId)
	if err != nil {
		return nil, err
	}

	stats = CollectionStatistics{
		ItemsCount:   numItems,
		OwnersCount:  numOwners,
		FloorPrice:   minPrice,
		VolumeTraded: sumPrice,
	}

	err = cache.GetCacher().Set(cacheKey, stats, StatisticsExpirePeriod)
	if err != nil {
		log.Debug("could not set cache", "err", err)
	}

	return &stats, nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getTokensRegisteredByUser(userAddress string, blockchainProxy string) ([]string, error) {
	var resp ProxyRegisteredNFTsResponse

	url := fmt.Sprintf(RegisteredNFTsBaseFormat, blockchainProxy, userAddress)
	err := cache.GetCacher().Get(url, &resp)
	if err == nil {
		return resp.Data.Tokens, nil
	}

	err = HttpGet(url, &resp)
	if err != nil {
		return nil, err
	}

	err = cache.GetCacher().Set(url, resp, HttpResponseExpirePeriod)
	if err != nil {
		log.Debug("could not cache response", "err", err)
	}

	return resp.Data.Tokens, nil
}

func checkValidInput(request *CreateCollectionRequest) error {
	if len(request.TokenId) == 0 {
		return errors.New("empty token id")
	}

	if len(request.TokenId) > TokenIdMaxLen {
		return errors.New("empty token id")
	}

	if len(request.Name) == 0 {
		return errors.New("empty name")
	}

	if len(request.Name) > MaxNameLen {
		return errors.New("name too long")
	}

	if len(request.Description) > MaxDescLen {
		return errors.New("description too long")
	}

	if len(request.Website) > MaxLinkLen {
		return errors.New("website too long")
	}

	if len(request.DiscordLink) > MaxLinkLen {
		return errors.New("discord link too long")
	}

	if len(request.TwitterLink) > MaxLinkLen {
		return errors.New("twitter link too long")
	}

	if len(request.InstagramLink) > MaxLinkLen {
		return errors.New("instagram link too long")
	}

	if len(request.TelegramLink) > MaxLinkLen {
		return errors.New("telegram link too long")
	}

	return nil
}
