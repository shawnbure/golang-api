package storage

import (
	"github.com/erdsea/erdsea-api/data"
	"gorm.io/gorm"
)

func AddNewAccount(account *data.Account) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&account)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetAccountById(id uint64) (*data.Account, error) {
	var account data.Account

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&account, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &account, nil
}

func GetAccountByAddress(name string) (*data.Account, error) {
	var account data.Account

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&account, "address = ?", name)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &account, nil
}
