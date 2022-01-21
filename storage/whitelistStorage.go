package storage

import (
	"gorm.io/gorm"

	"github.com/ENFT-DAO/youbei-api/data/entities"
)

func AddWhitelist(whitelist *entities.Whitelist) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&whitelist)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetWhitelistById(id uint64) (*entities.Whitelist, error) {
	var whitelist entities.Whitelist

	if id == 0 {
		return &whitelist, nil
	}

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&whitelist, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &whitelist, nil
}

func GetWhitelistByAddress(address string) (*entities.Whitelist, error) {
	var whitelist entities.Whitelist

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&whitelist, "address = ?", address)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &whitelist, nil
}

func GetWhitelistsByCollectionID(collectionID uint64) ([]entities.Whitelist, error) {
	var whitelists []entities.Whitelist

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&whitelists, "collectionID = ?", collectionID)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return whitelists, nil
}

func UpdateWhitelist(whitelist *entities.Whitelist) error {
	database, err := GetDBOrError()

	if err != nil {
		return err
	}

	txCreate := database.Save(&whitelist)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func UpdateWhitelistAmountWhereId(whitelistID uint64, amount uint64) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	tx := database.Table("whitelist").Where("ID = ?", whitelistID).Update("amount", amount)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
