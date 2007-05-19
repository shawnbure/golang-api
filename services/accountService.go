package services

import (
	"encoding/json"
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
	account := entities.Account{
		Address:       request.Address,
		Name:          request.Name,
		Description:   request.Description,
		Website:       request.Website,
		TwitterLink:   request.TwitterLink,
		InstagramLink: request.InstagramLink,
		CreatedAt:     uint64(time.Now().Unix()),
	}

	err := storage.AddAccount(&account)
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
		_, innerErr := tx.CreateBucketIfNotExists(WalletAddressToAccountCacheInfo)
		if innerErr != nil {
			return innerErr
		}

		innerErr = tx.Bucket(WalletAddressToAccountCacheInfo).Put([]byte(walletAddress), entryBytes)
		return innerErr
	})

	return &cacheInfo, nil
}

func GetAccountCacheInfo(walletAddress string) (*AccountCacheInfo, error) {
	db := cache.GetBolt()

	var bytes []byte
	err := db.View(func(tx *bolt.Tx) error {
		_, innerErr := tx.CreateBucketIfNotExists(WalletAddressToAccountCacheInfo)
		if innerErr != nil {
			return innerErr
		}

		bytes = tx.Bucket(WalletAddressToAccountCacheInfo).Get([]byte(walletAddress))
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
			return nil, err
		}

		cacheInfo, innerErr = AddAccountToCache(walletAddress, account.ID, account.Name)
		if innerErr != nil {
			return nil, err
		}
	}

	return cacheInfo, nil
}
