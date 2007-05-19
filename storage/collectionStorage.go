package storage

import (
	"github.com/erdsea/erdsea-api/data"
	"gorm.io/gorm"
)

func AddNewCollection(collection *data.Collection) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&collection)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetCollectionById(id uint64) (*data.Collection, error) {
	var collection data.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&collection, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &collection, nil
}

func GetCollectionsCreatedBy(id uint64) ([]data.Collection, error) {
	var collections []data.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&collections, "creator_id = ?", id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

func GetCollectionByName(name string) (*data.Collection, error) {
	var collection data.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&collection, "name = ?", name)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &collection, nil
}

func GetCollectionByTokenId(tokenId string) (*data.Collection, error) {
	var collection data.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&collection, "token_id = ?", tokenId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &collection, nil
}

func GetCollectionsWithOffsetLimit(offset int, limit int) ([]data.Collection, error) {
	var collections []data.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&collections)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

func GetCollectionsWithNameAlikeWithLimit(name string, limit int) ([]data.Collection, error) {
	var collections []data.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Limit(limit).Where("name LIKE ?", name).Find(&collections)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}
