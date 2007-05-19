package storage

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/erdsea/erdsea-api/data/entities"
)

func AddAsset(asset *entities.Asset) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&asset)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func UpdateAsset(asset *entities.Asset) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Save(&asset)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetAssetById(id uint64) (*entities.Asset, error) {
	var asset entities.Asset

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&asset, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &asset, nil
}

func GetAssetByTokenIdAndNonce(tokenId string, nonce uint64) (*entities.Asset, error) {
	var asset entities.Asset

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&asset, "token_id = ? AND nonce = ?", tokenId, nonce)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &asset, nil
}

func GetAssetsByOwnerIdWithOffsetLimit(ownerId uint64, offset int, limit int) ([]entities.Asset, error) {
	var assets []entities.Asset

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&assets, "owner_id = ?", ownerId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return assets, nil
}

func GetAssetsByCollectionId(collectionId uint64) ([]entities.Asset, error) {
	var assets []entities.Asset

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&assets, "collection_id = ?", collectionId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return assets, nil
}

func GetAssetsByCollectionIdWithOffsetLimit(collectionId uint64, offset int, limit int, attributesFilters map[string]string) ([]entities.Asset, error) {
	var assets []entities.Asset

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit)
	for k, v := range attributesFilters {
		txRead.Where(datatypes.JSONQuery("attributes").Equals(v, k))
	}

	txRead.Find(&assets, "collection_id = ?", collectionId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return assets, nil
}

func GetListedAssetsByCollectionIdWithOffsetLimit(collectionId uint64, offset int, limit int) ([]entities.Asset, error) {
	var assets []entities.Asset

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&assets, "listed = true AND collection_id = ?", collectionId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return assets, nil
}

func CountListedAssetsByCollectionId(collectionId uint64) (uint64, error) {
	count := int64(0)

	database, err := GetDBOrError()
	if err != nil {
		return 0, err
	}

	txRead := database.Model(&entities.Asset{}).Where("listed = true AND collection_id = ?", collectionId).Count(&count)
	if txRead.Error != nil {
		return 0, txRead.Error
	}

	return uint64(count), nil
}

func CountUniqueOwnersWithListedAssetsByCollectionId(collectionId uint64) (uint64, error) {
	count := int64(0)

	database, err := GetDBOrError()
	if err != nil {
		return 0, err
	}

	txRead := database.Model(&entities.Asset{}).Where("listed = true AND collection_id = ?", collectionId).Distinct("owner_id").Count(&count)
	if txRead.Error != nil {
		return 0, txRead.Error
	}

	return uint64(count), nil
}
