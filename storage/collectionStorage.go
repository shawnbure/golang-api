package storage

import (
	"github.com/erdsea/erdsea-api/data"
)

func AddNewCollection(collection *data.Collection) error {
	database := GetDB()
	if database == nil {
		return NoDBError
	}

	txCreate := GetDB().Create(&collection)
	if txCreate.Error != nil {
		return txCreate.Error
	}

	return nil
}

func GetCollectionById(id uint64) (*data.Collection, error) {
	var collection data.Collection

	database := GetDB()
	if database == nil {
		return nil, NoDBError
	}

	txRead := GetDB().Find(&collection, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &collection, nil
}

func GetCollectionsCreatedBy(id uint64) ([]data.Collection, error) {
	var collections []data.Collection

	database := GetDB()
	if database == nil {
		return nil, NoDBError
	}

	txRead := GetDB().Find(&collections, "creator_id = ?", id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

func GetCollectionByName(name string) (*data.Collection, error) {
	var collection data.Collection

	database := GetDB()
	if database == nil {
		return nil, NoDBError
	}

	txRead := GetDB().Find(&collection, "name = ?", name)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &collection, nil
}
