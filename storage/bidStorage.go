package storage

import (
	"gorm.io/gorm"

	"github.com/erdsea/erdsea-api/data/entities"
)

func AddBid(p *entities.Bid) error {
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

func DeleteBidsForTokenId(tokenDbId uint64) error {
	var bids []entities.Bid

	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Delete(bids, "token_id = ?", tokenDbId)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetBidsForTokenWithOffsetLimit(tokenId uint64, offset int, limit int) ([]entities.Bid, error) {
	var bids []entities.Bid

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&bids, "token_id = ?", tokenId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return bids, nil
}
