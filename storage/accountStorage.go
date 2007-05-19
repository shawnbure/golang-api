package storage

import (
	"gorm.io/gorm"

	"github.com/erdsea/erdsea-api/data/entities"
)

func AddAccount(account *entities.Account) error {
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

func UpdateAccount(account *entities.Account) error {
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

func UpdateAccountProfileWhereId(accountId uint64, link string) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	tx := database.Table("accounts").Where("id = ?", accountId).Update("profile_image_link", link)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func UpdateAccountCoverWhereId(accountId uint64, link string) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	tx := database.Table("accounts").Where("id = ?", accountId).Update("cover_image_link", link)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func GetAccountById(id uint64) (*entities.Account, error) {
	var account entities.Account

	if id == 0 {
		return &account, nil
	}

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

func GetAccountByAddress(address string) (*entities.Account, error) {
	var account entities.Account

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&account, "address = ?", address)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &account, nil
}

func GetAccountsWithNameAlikeWithLimit(name string, limit int) ([]entities.Account, error) {
	var accounts []entities.Account

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Limit(limit).Where("name ILIKE ?", name).Find(&accounts)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return accounts, nil
}
