package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data/entities"
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

type AccountCacheInfo struct {
	AccountId   uint64
	AccountName string
}

var (
	AccountSearchCacheKeyFormat     = "AccountSearch:%s"
	AccountSearchExpirePeriod       = 20 * time.Minute
	WalletAddressToAccountCacheInfo = []byte("walletToAcc")
)

func GetOrCreateAccount(address string) (*entities.Account, error) {
	account, err := storage.GetAccountByAddress(address)
	if err != nil {
		account = &entities.Account{
			Address:   address,
			CreatedAt: uint64(time.Now().Unix()),
		}

		err = storage.AddAccount(account)
		if err != nil {
			return nil, err
		}

		_, err = AddAccountToCache(account.Address, account.ID, account.Name)
		if err != nil {
			log.Debug("could not add account to cache")
		}
	}

	return account, nil
}

func CreateAccount(request *CreateAccountRequest) (*entities.Account, error) {
	err := checkValidCreateAccountRequest(request)
	if err != nil {
		return nil, err
	}

	account := entities.Account{
		Address:       request.Address,
		Name:          request.Name,
		Description:   request.Description,
		Website:       request.Website,
		TwitterLink:   request.TwitterLink,
		InstagramLink: request.InstagramLink,
		CreatedAt:     uint64(time.Now().Unix()),
	}

	err = storage.AddAccount(&account)
	if err != nil {
		return nil, err
	}

	_, err = AddAccountToCache(account.Address, account.ID, account.Name)
	if err != nil {
		log.Debug("could not add account to cache")
	}

	return &account, err
}

func UpdateAccount(account *entities.Account, request *SetAccountRequest) error {
	err := checkValidSetAccountRequest(request)
	if err != nil {
		return err
	}

	account.Name = request.Name
	account.Description = request.Description
	account.InstagramLink = request.InstagramLink
	account.TwitterLink = request.TwitterLink
	account.Website = request.Website

	return storage.UpdateAccount(account)
}

func GetAccountsWithNameAlike(name string, limit int) ([]entities.Account, error) {
	var byteArray []byte
	var accountArray []entities.Account

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

func AddAccountToCache(walletAddress string, accountId uint64, accountName string) (*AccountCacheInfo, error) {
	db := cache.GetBolt()
	cacheInfo := AccountCacheInfo{
		AccountId:   accountId,
		AccountName: accountName,
	}

	entryBytes, err := json.Marshal(&cacheInfo)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, innerErr := tx.CreateBucketIfNotExists(WalletAddressToAccountCacheInfo)
		if innerErr != nil {
			return innerErr
		}

		innerErr = bucket.Put([]byte(walletAddress), entryBytes)
		return innerErr
	})

	return &cacheInfo, nil
}

func GetAccountCacheInfo(walletAddress string) (*AccountCacheInfo, error) {
	db := cache.GetBolt()

	var bytes []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(WalletAddressToAccountCacheInfo)
		if bucket == nil {
			return errors.New("no bucket for account cache")
		}

		bytes = bucket.Get([]byte(walletAddress))
		return nil
	})
	if err != nil {
		return nil, err
	}

	var cacheInfo AccountCacheInfo
	err = json.Unmarshal(bytes, &cacheInfo)
	if err != nil {
		return nil, err
	}

	return &cacheInfo, nil
}

func GetOrAddAccountCacheInfo(walletAddress string) (*AccountCacheInfo, error) {
	cacheInfo, err := GetAccountCacheInfo(walletAddress)
	if err != nil {
		account, innerErr := storage.GetAccountByAddress(walletAddress)
		if innerErr != nil {
			return nil, innerErr
		}

		cacheInfo, innerErr = AddAccountToCache(walletAddress, account.ID, account.Name)
		if innerErr != nil {
			return nil, innerErr
		}
	}

	return cacheInfo, nil
}

func checkValidCreateAccountRequest(request *CreateAccountRequest) error {
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

	if len(request.TwitterLink) > MaxLinkLen {
		return errors.New("twitter link too long")
	}

	if len(request.InstagramLink) > MaxLinkLen {
		return errors.New("instagram link too long")
	}

	return nil
}

func checkValidSetAccountRequest(request *SetAccountRequest) error {
	if len(request.Description) > MaxDescLen {
		return errors.New("description too long")
	}

	if len(request.Website) > MaxLinkLen {
		return errors.New("website too long")
	}

	if len(request.TwitterLink) > MaxLinkLen {
		return errors.New("twitter link too long")
	}

	if len(request.InstagramLink) > MaxLinkLen {
		return errors.New("instagram link too long")
	}

	return nil
}
