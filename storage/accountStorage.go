package storage

import (
	"gorm.io/gorm"

	"github.com/erdsea/erdsea-api/data"
)

func AddAccount(account *data.Account) error {
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

func UpdateAccount(account *data.Account) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Save(&account)
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

func GetAccountsWithNameAlikeWithLimit(name string, limit int) ([]data.Account, error) {
	var accounts []data.Account

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Limit(limit).Where("name LIKE ?", name).Find(&accounts)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return accounts, nil
}
