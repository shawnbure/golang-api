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

func UpdateDeposit(deposit *entities.Deposit) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Save(&deposit)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
