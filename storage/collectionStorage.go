package storage

import (
	"github.com/erdsea/erdsea-api/data"
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

	return &collection, nil
}
