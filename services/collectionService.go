package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/stats/collstats"
	"github.com/erdsea/erdsea-api/storage"
)

const (
	TokenIdMaxLen                  = 15
	MaxNameLen                     = 20
	MaxLinkLen                     = 100
	MaxDescLen                     = 1000
	RegisteredNFTsBaseFormat       = "%s/address/%s/registered-nfts"
	HttpResponseExpirePeriod       = 10 * time.Minute
	CollectionSearchCacheKeyFormat = "CollectionSearch:%s"
	CollectionSearchExpirePeriod   = 20 * time.Minute
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

type UpdateCollectionRequest struct {
	Description   string `json:"description"`
	Website       string `json:"website"`
	DiscordLink   string `json:"discordLink"`
	TwitterLink   string `json:"twitterLink"`
	InstagramLink string `json:"instagramLink"`
	TelegramLink  string `json:"telegramLink"`
}

type ProxyRegisteredNFTsResponse struct {
	Data struct {
		Tokens []string `json:"tokens"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

func CreateCollection(request *CreateCollectionRequest, blockchainProxy string) (*entities.Collection, error) {
	err := checkValidInputOnCreate(request)
	if err != nil {
		return nil, err
	}

	tokensRegisteredByUser, err := getTokensRegisteredByUser(request.UserAddress, blockchainProxy)
	if err != nil {
		return nil, err
	}
	if !contains(tokensRegisteredByUser, request.TokenId) {
		return nil, errors.New("token not owner by user")
	}

	_, err = storage.GetCollectionByName(request.Name)
	if err == nil {
		return nil, errors.New("collection name already taken")
	}

	_, err = storage.GetCollectionByTokenId(request.TokenId)
	if err == nil {
		return nil, errors.New("token id already has a collection associated")
	}

	account, err := GetOrCreateAccount(request.UserAddress)
	if err != nil {
		return nil, err
	}

	collection := &entities.Collection{
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

	err = storage.AddCollection(collection)
	if err != nil {
		return nil, err
	}

	_, err = collstats.AddCollectionToCache(collection.ID, collection.Name, collection.TokenID)
	if err != nil {
		log.Debug("could not add to coll stats")
	}

	return collection, nil
}

func UpdateCollection(collection *entities.Collection, request *UpdateCollectionRequest) error {
	err := checkValidInputOnUpdate(request)

	collection.Description = request.Description
	collection.Website = request.Website
	collection.DiscordLink = request.DiscordLink
	collection.TwitterLink = request.TwitterLink
	collection.InstagramLink = request.InstagramLink
	collection.TelegramLink = request.TelegramLink

	err = storage.UpdateCollection(collection)
	if err != nil {
		return err
	}

	return nil
}

func GetCollectionsWithNameAlike(name string, limit int) ([]entities.Collection, error) {
	var byteArray []byte
	var collectionArray []entities.Collection

	cacheKey := fmt.Sprintf(CollectionSearchCacheKeyFormat, name)
	err := cache.GetCacher().Get(cacheKey, &byteArray)
	if err == nil {
		err = json.Unmarshal(byteArray, &collectionArray)
		return collectionArray, err
	}

	searchName := "%" + name + "%"
	collectionArray, err = storage.GetCollectionsWithNameAlikeWithLimit(searchName, limit)
	if err != nil {
		return nil, err
	}

	byteArray, err = json.Marshal(collectionArray)
	if err == nil {
		err = cache.GetCacher().Set(cacheKey, byteArray, CollectionSearchExpirePeriod)
		if err != nil {
			log.Debug("could not set cache", "err", err)
		}
	}

	return collectionArray, nil
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

func checkValidInputOnCreate(request *CreateCollectionRequest) error {
	if len(request.TokenId) == 0 {
		return errors.New("empty token id")
	}

	if len(request.TokenId) > TokenIdMaxLen {
		return errors.New("token id too long")
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

func checkValidInputOnUpdate(request *UpdateCollectionRequest) error {
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
