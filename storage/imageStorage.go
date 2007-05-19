package storage

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/erdsea/erdsea-api/data/entities/images"
)

func GetAccountProfileImageByAccountId(accountId uint64) (*images.AccountProfileImage, error) {
	var image images.AccountProfileImage

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&image, "account_id = ?", accountId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &image, nil
}

func AddOrUpdateAccountProfileImage(image *images.AccountProfileImage) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"image_base64"}),
	}).Create(image)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetAccountCoverImageByAccountId(accountId uint64) (*images.AccountCoverImage, error) {
	var image images.AccountCoverImage

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&image, "account_id = ?", accountId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &image, nil
}

func AddOrUpdateAccountCoverImage(image *images.AccountCoverImage) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"image_base64"}),
	}).Create(image)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetCollectionProfileImageByCollectionId(collectionId uint64) (*images.CollectionProfileImage, error) {
	var image images.CollectionProfileImage

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&image, "collection_id = ?", collectionId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &image, nil
}

func AddOrUpdateCollectionProfileImage(image *images.CollectionProfileImage) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "collection_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"image_base64"}),
	}).Create(image)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetCollectionCoverImageByCollectionId(collectionId uint64) (*images.CollectionCoverImage, error) {
	var image images.CollectionCoverImage

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&image, "collection_id = ?", collectionId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &image, nil
}

func AddOrUpdateCollectionCoverImage(image *images.CollectionCoverImage) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "collection_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"image_base64"}),
	}).Create(image)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
