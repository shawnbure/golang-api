package storage

import (
	"github.com/erdsea/erdsea-api/data/images"
	"gorm.io/gorm"
)

func GetAccountProfileImageByUserId(userId uint64) (*images.AccountProfileImage, error) {
	var image images.AccountProfileImage

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&image, "user_id = ?", userId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &image, nil
}

func AddNewAccountProfileImageByUserId(image *images.AccountProfileImage) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(image)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetAccountCoverImageByUserId(userId uint64) (*images.AccountCoverImage, error) {
	var image images.AccountCoverImage

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&image, "user_id = ?", userId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &image, nil
}

func AddNewAccountCoverImageByUserId(image *images.AccountCoverImage) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(image)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
