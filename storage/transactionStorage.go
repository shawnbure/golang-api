package storage

import (
	"github.com/erdsea/erdsea-api/data"
	"gorm.io/gorm"
)

func AddNewTransaction(transaction *data.Transaction) error {
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

func GetTransactionById(id uint64) (*data.Transaction, error) {
	var transaction data.Transaction

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

func GetTransactionsBySellerId(id uint64) ([]data.Transaction, error) {
	var transactions []data.Transaction

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

func GetTransactionsByBuyerId(id uint64) ([]data.Transaction, error) {
	var transactions []data.Transaction

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

func GetTransactionsByBuyerOrSellerId(id uint64) ([]data.Transaction, error) {
	var transactions []data.Transaction

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

func GetTransactionsByAssetId(id uint64) ([]data.Transaction, error) {
	var transactions []data.Transaction

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

func GetTransactionByHash(hash string) (*data.Transaction, error) {
	var transaction data.Transaction

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

