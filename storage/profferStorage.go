package storage

import (
	"gorm.io/gorm"

	"github.com/erdsea/erdsea-api/data/entities"
)

func AddProffer(p *entities.Proffer) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&p)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func DeleteProffersForTokenId(tokenDbId uint64) error {
	var proffers []entities.Proffer

	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Delete(proffers, "token_id = ?", tokenDbId)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func DeleteOfferByTokenIdAndAccountId(tokenDbId uint64, accountDbId uint64) error {
	var proffer entities.Proffer

	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Delete(&proffer, "type = ? AND token_id = ? AND offeror_id = ?", entities.Offer, tokenDbId, accountDbId)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetProffersForTokenWithOffsetLimit(tokenId uint64, offset int, limit int) ([]entities.Proffer, error) {
	var proffers []entities.Proffer

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&proffers, "token_id = ?", tokenId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return proffers, nil
}
