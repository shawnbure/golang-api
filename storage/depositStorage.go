package storage

import (
	"github.com/erdsea/erdsea-api/data/entities"
	"gorm.io/gorm"
)

func AddDeposit(deposit *entities.Deposit) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&deposit)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func UpdateDepositByOwnerId(deposit *entities.Deposit, ownerId uint64) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Table("deposits").Where("owner_id = ?", ownerId)
	txCreate.Update("amount_nominal", deposit.AmountNominal).Update("amount_string", deposit.AmountString)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
