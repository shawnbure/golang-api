package storage

import (
	"gorm.io/gorm"

	"github.com/erdsea/erdsea-api/data/entities"
)

func AddTransaction(transaction *entities.Transaction) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&transaction)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetTransactionById(id uint64) (*entities.Transaction, error) {
	var transaction entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&transaction, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &transaction, nil
}

func GetTransactionsBySellerId(id uint64) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&transactions, "seller_id = ?", id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetTransactionsByBuyerId(id uint64) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&transactions, "buyer_id = ?", id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetTransactionsByBuyerOrSellerId(id uint64) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&transactions, "seller_id = ? OR buyer_id = ?", id, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetTransactionsByBuyerOrSellerIdWithOffsetLimit(id uint64, offset int, limit int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&transactions, "seller_id = ? OR buyer_id = ?", id, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetTransactionsByAssetId(id uint64) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&transactions, "asset_id = ?", id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetTransactionsByAssetIdWithOffsetLimit(id uint64, offset int, limit int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&transactions, "asset_id = ?", id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetTransactionsByCollectionIdWithOffsetLimit(id uint64, offset int, limit int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&transactions, "collection_id = ?", id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetTransactionByHash(hash string) (*entities.Transaction, error) {
	var transaction entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&transaction, "hash = ?", hash)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &transaction, nil
}

func GetTransactionsWithOffsetLimit(offset int, limit int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&transactions)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetMinBuyPriceForTransactionsWithCollectionId(collectionId uint64) (float64, error) {
	var price float64

	database, err := GetDBOrError()
	if err != nil {
		return float64(0), err
	}

	txRead := database.Select("MIN(price_nominal)").
		Where("type = ? AND collection_id = ?", entities.BuyAsset, collectionId).
		Table("transactions").
		Find(&price)

	if txRead.Error != nil {
		return float64(0), txRead.Error
	}

	return price, nil
}

func GetSumBuyPriceForTransactionsWithCollectionId(collectionId uint64) (float64, error) {
	var price float64

	database, err := GetDBOrError()
	if err != nil {
		return float64(0), err
	}

	txRead := database.Select("SUM(price_nominal)").
		Where("type = ? AND collection_id = ?", entities.BuyAsset, collectionId).
		Table("transactions").
		Find(&price)

	if txRead.Error != nil {
		return float64(0), txRead.Error
	}

	return price, nil
}
