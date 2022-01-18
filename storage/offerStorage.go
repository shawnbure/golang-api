package storage

import (
	"gorm.io/gorm"

	"github.com/ENFT-DAO/youbei-api/data/entities"
)

func AddOffer(p *entities.Offer) error {
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

func DeleteOffersForTokenId(tokenDbId uint64) error {
	var offers []entities.Offer

	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Delete(offers, "token_id = ?", tokenDbId)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func DeleteOfferByOfferorForTokenId(offerorAddress string, tokenDbId uint64) error {
	var proffer entities.Offer

	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Delete(&proffer, "token_id = ? AND offeror_address = ?", tokenDbId, offerorAddress)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetOffersForTokenWithOffsetLimit(tokenId uint64, offset int, limit int) ([]entities.Offer, error) {
	var offer []entities.Offer

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Order("id desc").Find(&offer, "token_id = ?", tokenId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return offer, nil
}
