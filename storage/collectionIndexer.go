package storage

import (
	"github.com/ENFT-DAO/youbei-api/data/entities"
)

func GetCollectionIndexer(collectionAddr string) (entities.CollectionIndexer, error) {
	var stat entities.CollectionIndexer

	database, err := GetDBOrError()
	if err != nil {
		return stat, err
	}

	err = database.
		Model(&entities.CollectionIndexer{}).
		Where("collection_addr = ?", collectionAddr).
		Order("updated_at DESC").
		First(&stat).Error
	if err != nil {
		return stat, err
	}
	return stat, nil
}

func CreateCollectionStat(collectionAddr string) (stat entities.CollectionIndexer, err error) {
	database, err := GetDBOrError()
	if err != nil {
		return stat, err
	}
	stat.CollectionAddr = collectionAddr
	err = database.Create(&stat).Error
	if err != nil {
		return stat, err
	}
	return stat, nil
}

func UpdateCollectionIndexer(lastIndex uint64, collectionAddr string) (entities.CollectionIndexer, error) {
	var stat entities.CollectionIndexer

	database, err := GetDBOrError()
	if err != nil {
		return stat, err
	}

	stat.LastIndex = lastIndex
	stat.CollectionAddr = collectionAddr

	err = database.
		Model(&entities.CollectionIndexer{}).
		Where("collection_addr = ?", collectionAddr).
		Order("updated_at DESC").
		First(&stat).Error
	if err != nil {
		return stat, err
	}
	stat.LastIndex = lastIndex
	err = database.Save(&stat).Error
	// err = database.Updates(stat).
	// 	Where("id=?", stat.ID).Error
	if err != nil {
		return stat, err
	}
	return stat, nil
}
