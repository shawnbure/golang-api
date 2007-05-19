package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
)

type CreateAccountRequest struct {
	Address       string `json:"address"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Website       string `json:"website"`
	TwitterLink   string `json:"twitterLink"`
	InstagramLink string `json:"instagramLink"`
}

type SetAccountRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Website       string `json:"website"`
	TwitterLink   string `json:"twitterLink"`
	InstagramLink string `json:"instagramLink"`
}

var (
	AccountSearchCacheKeyFormat = "AccountSearch:%s"
	AccountSearchExpirePeriod   = 20 * time.Minute
)

func GetOrCreateAccount(address string) (*data.Account, error) {
	account, err := storage.GetAccountByAddress(address)
	if err != nil {
		account = &data.Account{
			Address:   address,
			CreatedAt: uint64(time.Now().Unix()),
		}

		err = storage.AddAccount(account)
		if err != nil {
			return nil, err
		}
	}

	return account, nil
}

func CreateAccount(request *CreateAccountRequest) (*data.Account, error) {
	account := data.Account{
		Address:       request.Address,
		Name:          request.Name,
		Description:   request.Description,
		Website:       request.Website,
		TwitterLink:   request.TwitterLink,
		InstagramLink: request.InstagramLink,
		CreatedAt:     uint64(time.Now().Unix()),
	}

	err := storage.AddAccount(&account)
	return &account, err
}

func UpdateAccount(account *data.Account, request *SetAccountRequest) error {
	account.Description = request.Description
	account.InstagramLink = request.InstagramLink
	account.TwitterLink = request.TwitterLink
	account.Website = request.Website
	account.Name = request.Name
	return storage.UpdateAccount(account)
}

func GetAccountsWithNameAlike(name string, limit int) ([]data.Account, error) {
	var byteArray []byte
	var accountArray []data.Account

	cacheKey := fmt.Sprintf(AccountSearchCacheKeyFormat, name)
	err := cache.GetCacher().Get(cacheKey, &byteArray)
	if err == nil {
		err = json.Unmarshal(byteArray, &accountArray)
		return accountArray, err
	}

	searchName := "%" + name + "%"
	accountArray, err = storage.GetAccountsWithNameAlikeWithLimit(searchName, limit)
	if err != nil {
		return nil, err
	}

	byteArray, err = json.Marshal(accountArray)
	if err == nil {
		err = cache.GetCacher().Set(cacheKey, byteArray, AccountSearchExpirePeriod)
		if err != nil {
			log.Debug("could not set cache", "err", err)
		}
	}

	return accountArray, nil
}
