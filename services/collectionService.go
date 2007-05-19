package services

import (
	"errors"
	"fmt"
	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
	"time"
)

const (
	TokenIdMaxLen            = 15
	MaxNameLen               = 20
	MaxLinkLen               = 100
	MaxDescLen               = 1000
	RegisteredNFTsBaseFormat = "%s/address/%s/registered-nfts"
	HttpResponseExpirePeriod = 10 * time.Minute
)

type CreateCollectionRequest struct {
	UserAddress   string `json:"UserAddress"`
	Name          string `json:"collectionName"`
	TokenId       string `json:"tokenId"`
	Description   string `json:"description"`
	Website       string `json:"website"`
	DiscordLink   string `json:"discordLink"`
	TwitterLink   string `json:"twitterLink"`
	InstagramLink string `json:"instagramLink"`
	TelegramLink  string `json:"telegramLink"`
}

type ElrondProxyRegisteredNFTsResponse struct {
	Data struct {
		Tokens []string `json:"tokens"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

func CreateCollection(request *CreateCollectionRequest, elrondProxy string) error {
	err := checkValidInput(request)
	if err != nil {
		return err
	}

	_, err = storage.GetCollectionByName(request.Name)
	if err == nil {
		return errors.New("collection name already taken")
	}

	tokenRegisteredByUser, err := getTokensRegisteredBy(request.UserAddress, elrondProxy)
	if err != nil {
		return err
	}
	if !contains(tokenRegisteredByUser, request.TokenId) {
		return errors.New("token not owner by user")
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
	}

	return storage.AddNewCollection(collection)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getTokensRegisteredBy(userAddress string, elrondProxy string) ([]string, error) {
	var resp ElrondProxyRegisteredNFTsResponse

	url := fmt.Sprintf(RegisteredNFTsBaseFormat, elrondProxy, userAddress)
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
		log.Debug("could not cache response")
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
