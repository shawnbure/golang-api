package storage

import (
	images2 "github.com/erdsea/erdsea-api/data/entities/images"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAccountProfileImageByAccountId(accountId uint64) (*images2.AccountProfileImage, error) {
	var image images2.AccountProfileImage

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

func AddOrUpdateAccountProfileImage(image *images2.AccountProfileImage) error {
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

func GetAccountCoverImageByAccountId(accountId uint64) (*images2.AccountCoverImage, error) {
	var image images2.AccountCoverImage

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

func AddOrUpdateAccountCoverImage(image *images2.AccountCoverImage) error {
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

func GetCollectionProfileImageByCollectionId(collectionId uint64) (*images2.CollectionProfileImage, error) {
	var image images2.CollectionProfileImage

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

func AddOrUpdateCollectionProfileImage(image *images2.CollectionProfileImage) error {
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

func GetCollectionCoverImageByCollectionId(collectionId uint64) (*images2.CollectionCoverImage, error) {
	var image images2.CollectionCoverImage

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

func AddOrUpdateCollectionCoverImage(image *images2.CollectionCoverImage) error {
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
