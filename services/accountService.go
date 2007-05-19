package services

import (
	"encoding/json"
	"fmt"
	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
	"time"
)

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

func AddOrUpdateAccount(account *data.Account) error {
	existingAccount, err := storage.GetAccountByAddress(account.Address)
	if err != nil {
		err = storage.AddNewAccount(account)
	} else {
		existingAccount.Description = account.Description
		existingAccount.InstagramLink = account.InstagramLink
		existingAccount.TwitterLink = account.TwitterLink
		existingAccount.Website = account.Website
		existingAccount.Name = account.Name
		err = storage.UpdateAccount(account)
	}

	return err
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
