package storage

import (
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"gorm.io/gorm"
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

func CreateCollectionStat(col entities.CollectionIndexer) (stat entities.CollectionIndexer, err error) {
	database, err := GetDBOrError()
	if err != nil {
		return stat, err
	}
	err = database.Create(&col).Error
	if err != nil {
		return stat, err
	}
	stat = col
	return stat, nil
}

func UpdateCollectionndexerWhere(collectionIndexer *entities.CollectionIndexer, toUpdate map[string]interface{}, where string, args ...interface{}) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Model(collectionIndexer).Where(where, args...).Updates(toUpdate)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
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
