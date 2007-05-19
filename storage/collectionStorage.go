package storage

import (
	"gorm.io/gorm"

	"github.com/erdsea/erdsea-api/data/entities"
)

func AddCollection(collection *entities.Collection) error {
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

func UpdateCollection(collection *entities.Collection) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Save(&collection)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func UpdateCollectionProfileWhereId(collectionId uint64, link string) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	tx := database.Table("collections").Where("id = ?", collectionId).Update("profile_image_link", link)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func UpdateCollectionCoverWhereId(collectionId uint64, link string) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	tx := database.Table("collections").Where("id = ?", collectionId).Update("cover_image_link", link)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func GetCollectionById(id uint64) (*entities.Collection, error) {
	var collection entities.Collection

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

func GetCollectionsCreatedBy(id uint64) ([]entities.Collection, error) {
	var collections []entities.Collection

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

func GetCollectionByName(name string) (*entities.Collection, error) {
	var collection entities.Collection

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

func GetCollectionWithNameILike(name string) (*entities.Collection, error) {
	var collection entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&collection, "name ILIKE ?", name)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &collection, nil
}

func GetCollectionByTokenId(tokenId string) (*entities.Collection, error) {
	var collection entities.Collection

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

func GetCollectionsWithOffsetLimit(offset int, limit int) ([]entities.Collection, error) {
	var collections []entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Order("priority desc").Find(&collections)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}

func GetCollectionsWithNameAlikeWithLimit(name string, limit int) ([]entities.Collection, error) {
	var collections []entities.Collection

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

func GetCollectionsByOwnerIdWithOffsetLimit(ownerId uint64, offset int, limit int) ([]entities.Collection, error) {
	var collections []entities.Collection

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&collections, "owner_id = ?", ownerId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return collections, nil
}
