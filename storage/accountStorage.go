package storage

import (
	"github.com/erdsea/erdsea-api/data"
)

func addNewAccount(account *data.Account) error {
	db := GetDB()
	if db == nil {
		return NoDBError
	}

	txCreate := GetDB().Create(&account)
	if txCreate.Error != nil {
		return txCreate.Error
	}

	return nil
}

func getAccountById(id uint64) (*data.Account, error) {
	var account data.Account

	db := GetDB()
	if db == nil {
		return nil, NoDBError
	}

	txRead := GetDB().Find(&account, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &account, nil
}

func getAccountByAddress(name string) (*data.Account, error) {
	var account data.Account

	db := GetDB()
	if db == nil {
		return nil, NoDBError
	}

	txRead := GetDB().Find(&account, "address = ?", name)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &account, nil
}
