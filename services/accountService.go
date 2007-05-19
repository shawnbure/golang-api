package services

import (
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

func GetOrCreateAccount(address string) (*data.Account, error) {
	account, err := storage.GetAccountByAddress(address)
	if err != nil {
		account = &data.Account{
			Address:   address,
			CreatedAt: uint64(time.Now().Unix()),
		}

		err = storage.AddNewAccount(account)
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

	return nil
}
