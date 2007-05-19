package storage

import (
	"errors"
	"github.com/erdsea/erdsea-api/data"
)

var NoDBError = errors.New("no DB Connection")

func addNewCollection(collection *data.Collection) error {
	db := GetDB()
	if db == nil {
		return NoDBError
	}

	txCreate := GetDB().Create(&collection)
	if txCreate.Error != nil {
		return txCreate.Error
	}

	return nil
}

func getCollectionById(id uint64) (*data.Collection, error) {
	var collection data.Collection

	db := GetDB()
	if db == nil {
		return nil, NoDBError
	}

	txRead := GetDB().Find(&collection, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &collection, nil
}

func getCollectionsCreatedBy(id uint64) ([]data.Collection, error) {
	var collections []data.Collection

	db := GetDB()
	if db == nil {
		return nil, NoDBError
	}

	txRead := GetDB().Find(&collections, "creator_id = ?", id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

func getCollectionByName(name string) (*data.Collection, error) {
	var collection data.Collection

	db := GetDB()
	if db == nil {
		return nil, NoDBError
	}

	txRead := GetDB().Find(&collection, "name = ?", name)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &collection, nil
}