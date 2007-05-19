package services

import (
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
)

func GetOrCreateAccount(address string) (*data.Account, error) {
	account, err := storage.GetAccountByAddress(address)
	if err != nil {
		account = &data.Account{
			Address: address,
		}

		err = storage.AddAccount(account)
		if err != nil {
			return nil, err
		}
	}

	return account, nil
}
