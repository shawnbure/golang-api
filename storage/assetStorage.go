package storage

import "github.com/erdsea/erdsea-api/data"

func AddNewAsset(asset *data.Asset) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&asset)
	if txCreate.Error != nil {
		return txCreate.Error
	}

	return nil
}

func GetAssetById(id uint64) (*data.Asset, error) {
	var asset data.Asset

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&asset, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &asset, nil
}

func GetAssetByTokenIdAndNonce(tokenId string, nonce uint64) (*data.Asset, error) {
	var asset data.Asset

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&asset, "token_id = ? AND nonce = ?", tokenId, nonce)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return &asset, nil
}

func GetAssetsOwnedBy(ownerId uint64) ([]data.Asset, error) {
	var assets []data.Asset

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&assets, "owner_id = ?", ownerId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return assets, nil
}

func GetAssetsByCollectionId(collectionId uint64) ([]data.Asset, error) {
	var assets []data.Asset

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
