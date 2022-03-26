package storage

import (
	"database/sql"

	"gorm.io/gorm"

	"github.com/ENFT-DAO/youbei-api/data/entities"
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

	txRead := database.Offset(offset).Limit(limit).Order("id desc").Find(&transactions, "seller_id = ? OR buyer_id = ?", id, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetTransactionsByTokenId(id uint64) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&transactions, "token_id = ?", id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetTransactionsByTokenIdWithOffsetLimit(id uint64, offset int, limit int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Order("timestamp desc").Find(&transactions, "token_id = ?", id)
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

	txRead := database.Offset(offset).Limit(limit).Order("id desc").Find(&transactions, "collection_id = ?", id)
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
func GetLastTokenTransaction(tokenId uint64) (entities.Transaction, error) {
	var transaction entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return transaction, err
	}

	txRead := database.Where("token_id=?", tokenId).Order("timestamp desc").Find(&transaction)
	if txRead.Error != nil {
		return transaction, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return transaction, gorm.ErrRecordNotFound
	}

	return transaction, nil
}
func GetTransactionWhere(where string, args ...interface{}) (entities.Transaction, error) {
	var transaction entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return transaction, err
	}

	txRead := database.Where(where, args...).Find(&transaction)
	if txRead.Error != nil {
		return transaction, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return transaction, gorm.ErrRecordNotFound
	}

	return transaction, nil
}

func GetTransactionsWithOffsetLimit(offset int, limit int) ([]entities.Transaction, error) {
	var transactions []entities.Transaction

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Order("id desc").Find(&transactions)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func DeleteTransaction(id uint64) error {
	var transaction entities.Transaction
	database, err := GetDBOrError()
	if err != nil {
		return err
	}
	txRead := database.Find(&transaction, "id = ?", id)
	if txRead.Error != nil {
		return txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func GetMinBuyPriceForTransactionsWithCollectionId(collectionId uint64) (float64, error) {
	var price float64

	database, err := GetDBOrError()
	if err != nil {
		return float64(0), err
	}

	nullFloat := sql.NullFloat64{}
	txRead := database.Select("MIN(price_nominal)").
		Where("type = ? AND collection_id = ?", entities.BuyToken, collectionId).
		Table("transactions").
		Find(&nullFloat)

	if txRead.Error != nil {
		return float64(0), txRead.Error
	}

	if nullFloat.Valid {
		price = nullFloat.Float64
	}

	return price, nil
}

func GetSumBuyPriceForTransactionsWithCollectionId(collectionId uint64) (float64, error) {
	var price float64

	database, err := GetDBOrError()
	if err != nil {
		return float64(0), err
	}

	nullFloat := sql.NullFloat64{}
	txRead := database.Select("SUM(price_nominal)").
		Where("type = ? AND collection_id = ?", entities.BuyToken, collectionId).
		Table("transactions").
		Find(&nullFloat)

	if txRead.Error != nil {
		return float64(0), txRead.Error
	}

	if nullFloat.Valid {
		price = nullFloat.Float64
	}

	return price, nil
}
