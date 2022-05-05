package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

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

func AddOrUpdateTransaction(transaction *entities.Transaction) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&transaction)
	if txCreate.Error != nil {
		if strings.Contains(txCreate.Error.Error(), "duplicate") {
			txCreate = database.Where("hash=?", transaction.Hash).Updates(transaction)
			if txCreate.Error != nil {
				return txCreate.Error
			}
		}
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
	tx := database.Where("id = ?", id).Delete(&transaction)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
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

func GetTransactionsCount() (int64, error) {
	var total int64

	database, err := GetDBOrError()
	if err != nil {
		return int64(0), err
	}

	txRead := database.
		Table("transactions").
		Count(&total)

	if txRead.Error != nil {
		return int64(0), txRead.Error
	}

	return total, nil
}

func GetTransactionsCountByType(_type entities.TxType) (int64, error) {
	var total int64

	database, err := GetDBOrError()
	if err != nil {
		return int64(0), err
	}

	txRead := database.
		Where("type = ?", _type).
		Table("transactions").
		Count(&total)

	if txRead.Error != nil {
		return int64(0), txRead.Error
	}

	return total, nil
}

func GetTransactionsCountByDate(date string) (int64, error) {
	var total int64

	database, err := GetDBOrError()
	if err != nil {
		return int64(0), err
	}

	txRead := database.
		Where("date_trunc('day', to_timestamp(transactions.timestamp))=?", date).
		Table("transactions").
		Count(&total)

	if txRead.Error != nil {
		return int64(0), txRead.Error
	}

	return total, nil
}

func GetTransactionsCountByDateAndType(_type entities.TxType, date string) (int64, error) {
	var total int64

	database, err := GetDBOrError()
	if err != nil {
		return int64(0), err
	}

	txRead := database.
		Where("date_trunc('day', to_timestamp(transactions.timestamp))=? and type=?", date, _type).
		Table("transactions").
		Count(&total)

	if txRead.Error != nil {
		return int64(0), txRead.Error
	}

	return total, nil
}

func DeleteAllTransaction() (int64, error) {
	database, err := GetDBOrError()
	if err != nil {
		return int64(0), err
	}

	tx := database.Where("1 = 1").Delete(&entities.Transaction{})
	if tx.Error != nil {
		return int64(0), err
	}

	return tx.RowsAffected, nil
}

func GetTotalTradedVolume() (*big.Float, error) {
	database, err := GetDBOrError()
	if err != nil {
		return big.NewFloat(0), err
	}

	var x sql.NullString

	txRead := database.
		Where("type=?", entities.BuyToken).
		Table("transactions").
		Select("sum(price_nominal)").
		Scan(&x)

	if txRead.Error != nil {
		return big.NewFloat(0), txRead.Error
	}

	if x.Valid {
		v, _ := new(big.Float).SetString(x.String)
		return v, nil
	}

	return big.NewFloat(0), errors.New("Null String ...")
}

func GetTotalTradedVolumeByDate(dateStr string) (*big.Float, error) {
	database, err := GetDBOrError()
	if err != nil {
		return big.NewFloat(0), err
	}

	var x sql.NullString

	txRead := database.
		Where("type=? and date_trunc('day', to_timestamp(transactions.timestamp))=?", entities.BuyToken, dateStr).
		Table("transactions").
		Select("sum(price_nominal)").
		Scan(&x)

	if txRead.Error != nil {
		return big.NewFloat(0), txRead.Error
	}

	if x.Valid {
		v, _ := new(big.Float).SetString(x.String)
		return v, nil
	}

	return big.NewFloat(0), errors.New("Null String ...")
}

func GetAllTransactionsWithPagination(lastTimestamp int64, currentPage, requestedPage, pageSize int, filter *entities.QueryFilter) ([]entities.TransactionDetail, error) {
	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	transactions := []entities.TransactionDetail{}

	query := ""
	order := "transactions.timestamp desc "
	offset := 0
	if lastTimestamp == 0 {
		query = filter.Query
	} else {
		query = "transactions.timestamp<?"
		if requestedPage < currentPage {
			query = "transactions.timestamp>?"
			order = "transactions.timestamp asc "
		}

		if requestedPage != currentPage {
			offset = (int(math.Abs(float64(requestedPage-currentPage))) - 1) * pageSize
		}

		if filter.Query != "" {
			query = fmt.Sprintf("(%s) and %s", filter.Query, query)
		}
		filter.Values = append(filter.Values, lastTimestamp)
	}

	txRead := database.Table("transactions").Select("transactions.type as tx_type, transactions.hash as tx_hash, transactions.id as tx_id, transactions.price_nominal as tx_price_nominal, transactions.timestamp as tx_timestamp, tokens.token_id as token_id, tokens.token_name as token_name, tokens.image_link as token_image_link, seller_account.address as from_address, transactions.buyer_id as to_id").
		Joins("inner join tokens on tokens.id=transactions.token_id ").
		Joins("inner join accounts as seller_account on seller_account.id=transactions.seller_id ").
		Order("transactions.timestamp desc ").
		Order(order).
		Where(query, filter.Values...).
		Offset(offset).
		Limit(pageSize).
		Scan(&transactions)

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetLast24HoursSalesTransactions(fromTime string, toTime string) ([]entities.TransactionDetail, error) {
	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	transactions := []entities.TransactionDetail{}
	txRead := database.Table("transactions").Select("transactions.type as tx_type, transactions.hash as tx_hash, transactions.id as tx_id, transactions.price_nominal as tx_price_nominal, transactions.timestamp as tx_timestamp, tokens.token_id as token_id, tokens.token_name as token_name, tokens.image_link as token_image_link, seller_account.address as from_address, buyer_account.address as to_address").
		Joins("inner join tokens on tokens.id=transactions.token_id ").
		Joins("inner join accounts as seller_account on seller_account.id=transactions.seller_id ").
		Joins("inner join accounts as buyer_account on buyer_account.id=transactions.buyer_id ").
		Order("transactions.timestamp desc").
		Where("date_trunc('hour', to_timestamp(transactions.timestamp))<? and date_trunc('hour', to_timestamp(transactions.timestamp))>=? and transactions.type=?", toTime, fromTime, entities.BuyToken).
		Scan(&transactions)

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetLast24HoursTotalVolume(fromTime, toTime string) (*big.Float, error) {
	database, err := GetDBOrError()
	if err != nil {
		return big.NewFloat(0), err
	}

	var x sql.NullString

	txRead := database.
		Where("date_trunc('hour', to_timestamp(transactions.timestamp))<? and date_trunc('hour', to_timestamp(transactions.timestamp))>=? and transactions.type=?", toTime, fromTime, entities.BuyToken).
		Table("transactions").
		Select("sum(price_nominal)").
		Scan(&x)

	if txRead.Error != nil {
		return big.NewFloat(0), txRead.Error
	}

	if x.Valid {
		v, _ := new(big.Float).SetString(x.String)
		return v, nil
	}

	return big.NewFloat(0), errors.New("Null String ...")
}

func GetAllActivitiesWithPagination(lastTimestamp int64, currentPage, requestedPage, pageSize int, filter *entities.QueryFilter) ([]entities.Activity, error) {
	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	transactions := []entities.Activity{}

	query := ""
	order := "transactions.timestamp desc "
	offset := 0
	if lastTimestamp == 0 {
		query = filter.Query
	} else {
		query = "transactions.timestamp<?"
		if requestedPage < currentPage {
			query = "transactions.timestamp>?"
			order = "transactions.timestamp asc "
		}

		if requestedPage != currentPage {
			offset = (int(math.Abs(float64(requestedPage-currentPage))) - 1) * pageSize
		}

		if filter.Query != "" {
			query = fmt.Sprintf("(%s) and %s", filter.Query, query)
		}
		filter.Values = append(filter.Values, lastTimestamp)
	}

	txRead := database.Table("transactions").Select("transactions.type as tx_type, transactions.hash as tx_hash, transactions.id as tx_id, transactions.price_nominal as tx_price_nominal, transactions.timestamp as tx_timestamp, tokens.token_id as token_id, tokens.token_name as token_name, tokens.image_link as token_image_link, seller_account.address as from_address, transactions.buyer_id as to_id, collections.id as collection_id, collections.name as collection_name, collections.token_id as collection_token_id").
		Joins("inner join tokens on tokens.id=transactions.token_id ").
		Joins("inner join collections on collections.id=transactions.collection_id ").
		Joins("inner join accounts as seller_account on seller_account.id=transactions.seller_id ").
		Order(order).
		Where(query, filter.Values...).
		Offset(offset).
		Limit(pageSize).
		Scan(&transactions)

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return transactions, nil
}

func GetTopBestSellerLastWeek(limit int, fromDateTimestamp string, toDateTimestamp string) ([]entities.TopVolumeByAddress, error) {
	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	records := []entities.TopVolumeByAddress{}
	txRead := database.Table("transactions").
		Select("seller_account.address as address, sum(transactions.price_nominal) as volume").
		Joins("inner join accounts as seller_account on seller_account.id=transactions.seller_id ").
		Group("address").
		Where("date_trunc('day', to_timestamp(transactions.timestamp))>=? and date_trunc('day', to_timestamp(transactions.timestamp))<? and transactions.type=?", fromDateTimestamp, toDateTimestamp, entities.BuyToken).
		Limit(limit).Order("volume desc").
		Scan(&records)

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return records, nil
}

func GetTopBestSellerLastWeekTransactions(fromDateTimestamp string, toDateTimestamp string, addresses []string) ([]entities.TransactionDetail, error) {
	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	records := []entities.TransactionDetail{}
	txRead := database.Table("transactions").
		Select("transactions.type as tx_type, transactions.hash as tx_hash, transactions.id as tx_id, transactions.price_nominal as tx_price_nominal, transactions.timestamp as tx_timestamp, tokens.token_id as token_id, tokens.token_name as token_name, tokens.image_link as token_image_link, seller_account.address as from_address, buyer_account.address as to_address").
		Joins("inner join tokens on tokens.id=transactions.token_id ").
		Joins("inner join accounts as seller_account on seller_account.id=transactions.seller_id ").
		Joins("inner join accounts as buyer_account on buyer_account.id=transactions.buyer_id ").
		Where("date_trunc('day', to_timestamp(transactions.timestamp))>=? and date_trunc('day', to_timestamp(transactions.timestamp))<? and seller_account.address in (?) and transactions.type=?", fromDateTimestamp, toDateTimestamp, addresses, entities.BuyToken).
		Order("from_address asc").
		Scan(&records)

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return records, nil
}

func GetTopBestBuyerLastWeek(limit int, fromDateTimestamp string, toDateTimestamp string) ([]entities.TopVolumeByAddress, error) {
	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	records := []entities.TopVolumeByAddress{}
	txRead := database.Table("transactions").
		Select("buyer_account.address as address, sum(transactions.price_nominal) as volume").
		Joins("inner join accounts as buyer_account on buyer_account.id=transactions.buyer_id ").
		Group("address").
		Where("date_trunc('day', to_timestamp(transactions.timestamp))>=? and date_trunc('day', to_timestamp(transactions.timestamp))<? and transactions.type=?", fromDateTimestamp, toDateTimestamp, entities.BuyToken).
		Limit(limit).Order("volume desc").
		Scan(&records)

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return records, nil
}

func GetTopBestBuyerLastWeekTransactions(fromDateTimestamp string, toDateTimestamp string, addresses []string) ([]entities.TransactionDetail, error) {
	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	records := []entities.TransactionDetail{}
	txRead := database.Table("transactions").
		Select("transactions.type as tx_type, transactions.hash as tx_hash, transactions.id as tx_id, transactions.price_nominal as tx_price_nominal, transactions.timestamp as tx_timestamp, tokens.token_id as token_id, tokens.token_name as token_name, tokens.image_link as token_image_link, seller_account.address as from_address, buyer_account.address as to_address").
		Joins("inner join tokens on tokens.id=transactions.token_id ").
		Joins("inner join accounts as seller_account on seller_account.id=transactions.seller_id ").
		Joins("inner join accounts as buyer_account on buyer_account.id=transactions.buyer_id ").
		Where("date_trunc('day', to_timestamp(transactions.timestamp))>=? and date_trunc('day', to_timestamp(transactions.timestamp))<? and buyer_account.address in (?) and transactions.type=?", fromDateTimestamp, toDateTimestamp, addresses, entities.BuyToken).
		Order("to_address asc").
		Scan(&records)

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return records, nil
}

func GetLast24HoursVerifiedListingTransactions(fromDateTimestamp string, toDateTimestamp string) ([]entities.VerifiedListingTransaction, error) {
	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	records := []entities.VerifiedListingTransaction{}
	txRead := database.Table("transactions").
		Select("transactions.hash as tx_hash, transactions.id as tx_id, transactions.price_nominal as tx_price_nominal, transactions.timestamp as tx_timestamp, tokens.token_name as token_name, tokens.image_link as token_image_link, seller_account.address as address, collections.name as collection_name, collections.token_id as collection_token_id").
		Joins("inner join tokens on tokens.id=transactions.token_id ").
		Joins("inner join accounts as seller_account on seller_account.id=transactions.seller_id ").
		Joins("inner join collections on collections.id=transactions.collection_id").
		Where("date_trunc('hour', to_timestamp(transactions.timestamp))>=? and date_trunc('hour', to_timestamp(transactions.timestamp))<? and collections.is_verified=? and transactions.type=?", fromDateTimestamp, toDateTimestamp, true, entities.ListToken).
		Scan(&records)

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return records, nil
}
