package storage

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/stretchr/testify/require"
)

func Test_AddTransaction(t *testing.T) {
	connectToTestDb()

	transaction := defaultTransaction()
	err := AddTransaction(&transaction)
	require.Nil(t, err)

	var transactionRead entities.Transaction
	txRead := GetDB().Last(&transactionRead)

	require.Nil(t, txRead.Error)
	require.Equal(t, transactionRead, transaction)
}

func Test_GetTransactionById(t *testing.T) {
	connectToTestDb()

	transaction := defaultTransaction()
	err := AddTransaction(&transaction)
	require.Nil(t, err)

	transactionRead, err := GetTransactionById(transaction.ID)
	require.Nil(t, err)
	require.Equal(t, transactionRead, &transaction)
}

func Test_GetTransactionsByBuyerId(t *testing.T) {
	connectToTestDb()

	transaction := defaultTransaction()
	err := AddTransaction(&transaction)
	require.Nil(t, err)

	otherTransaction := defaultTransaction()
	err = AddTransaction(&otherTransaction)
	require.Nil(t, err)

	transactionsRead, err := GetTransactionsByBuyerId(transaction.BuyerID)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(transactionsRead), 2)

	for _, transactionRead := range transactionsRead {
		require.Equal(t, transactionRead.BuyerID, transaction.BuyerID)
	}
}

func Test_GetTransactionsBySellerId(t *testing.T) {
	connectToTestDb()

	transaction := defaultTransaction()
	err := AddTransaction(&transaction)
	require.Nil(t, err)

	otherTransaction := defaultTransaction()
	err = AddTransaction(&otherTransaction)
	require.Nil(t, err)

	transactionsRead, err := GetTransactionsBySellerId(transaction.SellerID)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(transactionsRead), 2)

	for _, transactionRead := range transactionsRead {
		require.Equal(t, transactionRead.SellerID, transaction.SellerID)
	}
}

func Test_GetTransactionsByBuyerOrSellerId(t *testing.T) {
	connectToTestDb()

	transaction := defaultTransaction()
	err := AddTransaction(&transaction)
	require.Nil(t, err)

	otherTransaction := defaultTransaction()
	err = AddTransaction(&otherTransaction)
	require.Nil(t, err)

	transactionsRead, err := GetTransactionsByBuyerOrSellerId(transaction.SellerID)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(transactionsRead), 2)

	for _, transactionRead := range transactionsRead {
		sameSeller := transactionRead.SellerID == transaction.SellerID
		sameBuyer := transactionRead.BuyerID == transaction.BuyerID
		require.Equal(t, sameBuyer || sameSeller, true)
	}
}

func Test_GetTransactionsByTokenId(t *testing.T) {
	connectToTestDb()

	transaction := defaultTransaction()
	err := AddTransaction(&transaction)
	require.Nil(t, err)

	otherTransaction := defaultTransaction()
	err = AddTransaction(&otherTransaction)
	require.Nil(t, err)

	transactionsRead, err := GetTransactionsByTokenId(transaction.TokenID)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(transactionsRead), 2)

	for _, transactionRead := range transactionsRead {
		require.Equal(t, transactionRead.TokenID, transaction.TokenID)
	}
}

func Test_GetTransactionsByHash(t *testing.T) {
	connectToTestDb()

	transaction := defaultTransaction()
	transaction.Hash = "my_unique_hash"
	err := AddTransaction(&transaction)
	require.Nil(t, err)

	transactionRead, err := GetTransactionByHash(transaction.Hash)
	require.Nil(t, err)
	require.Equal(t, transactionRead.Hash, transaction.Hash)
}

func Test_GetTransactionWithMinPriceByCollectionId(t *testing.T) {
	connectToTestDb()

	transaction := defaultTransaction()
	transaction.PriceNominal = float64(1)
	transaction.Type = "Buy"
	transaction.Hash = "my_unique_hash"
	err := AddTransaction(&transaction)
	require.Nil(t, err)

	minPrice, err := GetMinBuyPriceForTransactionsWithCollectionId(99)
	require.Nil(t, err)
	require.Equal(t, minPrice, float64(1))
}

func Test_GetSumBuyPriceForTransactionsWithCollectionId(t *testing.T) {
	connectToTestDb()

	transaction := defaultTransaction()
	transaction.Type = "Buy"
	transaction.Hash = "my_unique_hash"
	err := AddTransaction(&transaction)
	require.Nil(t, err)

	sumPrice, err := GetSumBuyPriceForTransactionsWithCollectionId(1)
	require.Nil(t, err)
	require.GreaterOrEqual(t, sumPrice, float64(1_000_000_000_000_000_000_000))
}

func Test_GetTotalTradesCount(t *testing.T) {
	connectToTestDb()

	// first clean all transactions
	err := cleanTransactionTable()
	require.Nil(t, err)

	// insert some records
	err = insertBulkTransactions()
	require.Nil(t, err)

	t.Run("Get Total Trades", func(t *testing.T) {
		total, err := GetTransactionsCount()
		require.Nil(t, err)
		require.Equal(t, int64(4), total, "Two result does not match")
	})

	t.Run("Get Total Trades By Type", func(t *testing.T) {
		total, err := GetTransactionsCountByType(entities.BuyToken)
		require.Nil(t, err)
		require.Equal(t, int64(1), total, "Total buy trades does not match")
	})

	t.Run("Get Total Trades By Date", func(t *testing.T) {
		total, err := GetTransactionsCountByDate("2020-04-10")
		require.Nil(t, err)
		require.Equal(t, int64(2), total, "Total trades does not match")
	})

	t.Run("Get Total Trades By Date And Type", func(t *testing.T) {
		total, err := GetTransactionsCountByDateAndType(entities.BuyToken, "2020-04-10")
		require.Nil(t, err)
		require.Equal(t, int64(0), total, "Total buy trades does not match")
	})
}

func Test_GetTotalTradedVolume(t *testing.T) {
	connectToTestDb()

	// first clean all transactions
	err := cleanTransactionTable()
	require.Nil(t, err)

	//insert some records
	err = insertBulkTransactions()
	require.Nil(t, err)

	t.Run("Get Total Volumes Traded all the time", func(t *testing.T) {
		total, err := GetTotalTradedVolume()
		require.Nil(t, err)
		v, _ := new(big.Int).SetString("2000000000000000000000", 10)
		require.Equal(t, v, total, "Total traded volume does not match")
	})

	t.Run("Get Total For two different days", func(t *testing.T) {
		total, err := GetTotalTradedVolumeByDate("2020-04-09")
		require.Nil(t, err)
		v, _ := new(big.Int).SetString("1000000000000000000000", 10)
		require.Equal(t, v, total, "Total traded volume for specific date does not match")

		total, err = GetTotalTradedVolumeByDate("2020-04-10")
		require.Nil(t, err)
		v, _ = new(big.Int).SetString("1000000000000000000000", 10)
		require.Equal(t, v, total, "Total traded volume for specific date does not match")
	})
}

func Test_GetAllTransactionsWithDetail(t *testing.T) {
	connectToTestDb()

	// first clean all transactions
	//err := cleanTransactionTable()
	//require.Nil(t, err)

	//insert some records
	//err = insertBulkTransactions()
	//require.Nil(t, err)

	t.Run("Get Transactions With Detail", func(t *testing.T) {
		lastTimestamp := int64(-1)
		howMuchRow := 2
		filter := entities.QueryFilter{}

		transactions, err := GetAllTransactionsWithPagination(lastTimestamp, 1, 1, howMuchRow, &filter)
		require.Nil(t, err)
		require.Equal(t, len(transactions), 2, "The returned transactions array length does not matched")
	})

	t.Run("Get Transactions With Detail with pagination", func(t *testing.T) {
		lastTimestamp := int64(1586480400)
		howMuchRow := 2
		filter := entities.QueryFilter{}

		transactions, err := GetAllTransactionsWithPagination(lastTimestamp, 1, 1, howMuchRow, &filter)
		require.Nil(t, err)
		require.Equal(t, len(transactions), 2, "The returned transactions array length does not matched")
	})
}

func Test_GetDailySales(t *testing.T) {
	connectToTestDb()

	// first clean all transactions
	err := cleanTransactionTable()
	require.Nil(t, err)

	//insert some records
	err = insertBulkTransactions()
	require.Nil(t, err)

	t.Run("Get the last 24 hours transactions", func(t *testing.T) {
		fromTime := "2020-04-23 12:00:00"
		toTime := "2020-04-24 12:00:00"

		transactions, err := GetLast24HoursSalesTransactions(fromTime, toTime)
		require.Nil(t, err)

		require.Equal(t, len(transactions), 2, "The returned transactions count is incorrect")
	})

	t.Run("Get the total volume last 24 hours", func(t *testing.T) {
		fromTime := "2020-04-23 12:00:00"
		toTime := "2020-04-24 12:00:00"

		total, err := GetLast24HoursTotalVolume(fromTime, toTime)
		require.Nil(t, err)
		v, _ := new(big.Float).SetString("2000000000000000000000")
		require.Equal(t, total, v, "The total volume is not correct")
	})
}

func Test_GetAllActivities(t *testing.T) {
	connectToTestDb()

	// first clean all transactions
	//err := cleanTransactionTable()
	//require.Nil(t, err)

	//insert some records
	//err = insertBulkTransactions()
	//require.Nil(t, err)

	t.Run("Get all activities and check the list", func(t *testing.T) {
		lastTimestamp := int64(0)
		howMuchRow := 3
		filter := entities.QueryFilter{}
		collectFilter := entities.QueryFilter{}

		transactions, err := GetAllActivitiesWithPagination(lastTimestamp, 1, 1, howMuchRow, &filter, &collectFilter)
		require.Nil(t, err)

		require.Equal(t, len(transactions), 3, "The returned transactions array length does not matched")
	})

	t.Run("Get all activities With Detail with pagination", func(t *testing.T) {
		lastTimestamp := int64(1586480400)
		howMuchRow := 2
		filter := entities.QueryFilter{}
		collectFilter := entities.QueryFilter{}

		transactions, err := GetAllActivitiesWithPagination(lastTimestamp, 1, 1, howMuchRow, &filter, &collectFilter)
		require.Nil(t, err)

		require.Equal(t, len(transactions), 0, "The returned transactions array length does not matched")
	})

	t.Run("Get filtered activities", func(t *testing.T) {
		lastTimestamp := int64(0)
		howMuchRow := 10
		filter := entities.QueryFilter{
			Query:  "transactions.type=? OR transactions.type=?",
			Values: []interface{}{"List", "Buy"},
		}

		collectFilter := entities.QueryFilter{}

		transactions, err := GetAllActivitiesWithPagination(lastTimestamp, 0, 0, howMuchRow, &filter, &collectFilter)
		require.Nil(t, err)

		require.Equal(t, len(transactions), 2, "The returned transactions array length does not matched")
	})

}

func Test_GetWeeklyReport(t *testing.T) {
	connectToTestDb()

	//first clean all transactions
	err := cleanTransactionTable()
	require.Nil(t, err)

	//insert some records
	err = insertBulkTransactions()
	require.Nil(t, err)

	t.Run("Get best seller per week", func(t *testing.T) {
		howMuch := 10
		fromDate := "2020-04-22"
		toDate := "2020-04-27"

		records, err := GetTopBestSellerLastWeek(howMuch, fromDate, toDate)
		require.Nil(t, err)
		require.Equal(t, len(records), 2, "The result does not match")

		for _, r := range records {
			if r.Volume != float64(2_000_000_000_000_000_000_000) && r.Volume != float64(3_000_000_000_000_000_000_000) {
				require.Error(t, errors.New("The volumes do not matched properly"))
			}
		}
	})

	t.Run("Test Limit of returned result", func(t *testing.T) {
		howMuch := 1
		fromDate := "2020-04-22"
		toDate := "2020-04-27"

		records, err := GetTopBestSellerLastWeek(howMuch, fromDate, toDate)
		require.Nil(t, err)
		require.Equal(t, len(records), 1, "The result does not match")
	})

	t.Run("Get Transactions of best seller for last week", func(t *testing.T) {
		address1 := "erd123"
		address2 := "erd1234"
		fromDate := "2020-04-22"
		toDate := "2020-04-27"

		records, err := GetTopBestSellerLastWeekTransactions(fromDate, toDate, []string{address1, address2})
		require.Nil(t, err)
		require.Equal(t, len(records), 4, "The result does not match")
	})

	t.Run("Get best buyers per week", func(t *testing.T) {
		howMuch := 10
		fromDate := "2020-04-22"
		toDate := "2020-04-27"

		records, err := GetTopBestBuyerLastWeek(howMuch, fromDate, toDate)
		require.Nil(t, err)
		require.Equal(t, len(records), 2, "The result does not match")

		for _, r := range records {
			if r.Volume != float64(2_000_000_000_000_000_000_000) && r.Volume != float64(3_000_000_000_000_000_000_000) {
				require.Error(t, errors.New("The volumes do not matched properly"))
			}
		}
	})

	t.Run("Get Transactions of best buyers for last week", func(t *testing.T) {
		address1 := "erd123"
		address2 := "erd1234"
		fromDate := "2020-04-22"
		toDate := "2020-04-27"

		records, err := GetTopBestBuyerLastWeekTransactions(fromDate, toDate, []string{address1, address2})
		require.Nil(t, err)
		require.Equal(t, len(records), 4, "The result does not match")
	})
}

func Test_DailyReportOfListingTransactions(t *testing.T) {
	connectToTestDb()

	//first clean all transactions
	err := cleanTransactionTable()
	require.Nil(t, err)

	//insert some records
	err = insertBulkTransactions()
	require.Nil(t, err)

	t.Run("Daily Report of Verified Transactions of Type=List", func(t *testing.T) {
		fromTime := "2020-04-23 12:00:00"
		toTime := "2020-04-24 12:00:00"

		records, err := GetLast24HoursVerifiedListingTransactions(fromTime, toTime)
		require.Nil(t, err)

		require.Equal(t, len(records), 1, "The returned result is empty")
	})
}

func Test_hourlyAggregatedTrades(t *testing.T) {
	connectToTestDb()

	//first clean all transactions
	err := cleanTransactionTable()
	require.Nil(t, err)

	//insert some records
	err = insertBulkTransactions2()
	require.Nil(t, err)

	t.Run("Hourly Traded Volume", func(t *testing.T) {
		fromTime := "2020-04-23 20:00:00"
		toTime := "2020-04-23 21:00:00"

		result, err := GetAggregatedTradedVolumeHourly(fromTime, toTime, entities.BuyToken)
		require.Nil(t, err)
		require.Equal(t, result, float64(2_000_000_000_000_000_000_000), "Buy Volume don't match")

		result, err = GetAggregatedTradedVolumeHourly(fromTime, toTime, entities.ListToken)
		require.Nil(t, err)
		require.Equal(t, result, float64(2_000_000_000_000_000_000_000), "List Volume don't match")

		result, err = GetAggregatedTradedVolumeHourly(fromTime, toTime, entities.WithdrawToken)
		require.Nil(t, err)
		require.Equal(t, result, float64(1_000_000_000_000_000_000_000), "Withdraw Volume don't match")
	})
}

func defaultTransaction() entities.Transaction {
	return entities.Transaction{
		Hash:         "hash",
		Type:         "test",
		PriceNominal: 1_000_000_000_000_000_000_000,
		SellerID:     1,
		BuyerID:      2,
		TokenID:      5,
		CollectionID: 1,
	}
}

func cleanTransactionTable() error {
	_, err := DeleteAllTransaction()
	if err != nil {
		return err
	}

	return nil
}

func insertBulkTransactions() error {
	transaction := defaultTransaction()
	transaction.Type = entities.BuyToken
	transaction.Hash = "my_unique_hash1"
	transaction.Timestamp = uint64(time.Date(2020, 04, 23, 20, 2, 0, 0, time.UTC).Unix())
	err := AddTransaction(&transaction)
	if err != nil {
		return err
	}

	transaction = defaultTransaction()
	transaction.Type = entities.BuyToken
	transaction.Hash = "my_unique_hash5"
	transaction.SellerID = 3
	transaction.BuyerID = 1
	transaction.Timestamp = uint64(time.Date(2020, 04, 24, 1, 0, 0, 0, time.UTC).Unix())
	err = AddTransaction(&transaction)
	if err != nil {
		return err
	}

	transaction = defaultTransaction()
	transaction.Type = entities.WithdrawToken
	transaction.Hash = "my_unique_hash2"
	transaction.Timestamp = uint64(time.Date(2020, 04, 23, 15, 0, 0, 0, time.UTC).Unix())
	err = AddTransaction(&transaction)
	if err != nil {
		return err
	}

	transaction = defaultTransaction()
	transaction.Type = entities.AuctionToken
	transaction.Hash = "my_unique_hash3"
	transaction.SellerID = 3
	transaction.BuyerID = 1
	transaction.Timestamp = uint64(time.Date(2020, 04, 24, 12, 3, 0, 0, time.UTC).Unix())
	err = AddTransaction(&transaction)
	if err != nil {
		return err
	}

	transaction = defaultTransaction()
	transaction.Type = entities.ListToken
	transaction.Hash = "my_unique_hash4"
	transaction.Timestamp = uint64(time.Date(2020, 04, 10, 1, 0, 0, 0, time.UTC).Unix())
	err = AddTransaction(&transaction)
	if err != nil {
		return err
	}

	return nil
}

func insertBulkTransactions2() error {
	transaction := defaultTransaction()
	transaction.Type = entities.BuyToken
	transaction.Hash = "my_unique_hash1"
	transaction.Timestamp = uint64(time.Date(2020, 04, 23, 20, 2, 0, 0, time.UTC).Unix())
	err := AddTransaction(&transaction)
	if err != nil {
		return err
	}

	transaction = defaultTransaction()
	transaction.Type = entities.BuyToken
	transaction.Hash = "my_unique_hash5"
	transaction.SellerID = 3
	transaction.BuyerID = 1
	transaction.Timestamp = uint64(time.Date(2020, 04, 23, 20, 0, 0, 0, time.UTC).Unix())
	err = AddTransaction(&transaction)
	if err != nil {
		return err
	}

	transaction = defaultTransaction()
	transaction.Type = entities.WithdrawToken
	transaction.Hash = "my_unique_hash2"
	transaction.Timestamp = uint64(time.Date(2020, 04, 23, 20, 40, 0, 0, time.UTC).Unix())
	err = AddTransaction(&transaction)
	if err != nil {
		return err
	}

	transaction = defaultTransaction()
	transaction.Type = entities.ListToken
	transaction.Hash = "my_unique_hash3"
	transaction.SellerID = 3
	transaction.BuyerID = 1
	transaction.Timestamp = uint64(time.Date(2020, 04, 23, 20, 59, 0, 0, time.UTC).Unix())
	err = AddTransaction(&transaction)
	if err != nil {
		return err
	}

	transaction = defaultTransaction()
	transaction.Type = entities.ListToken
	transaction.Hash = "my_unique_hash4"
	transaction.Timestamp = uint64(time.Date(2020, 04, 23, 20, 3, 0, 0, time.UTC).Unix())
	err = AddTransaction(&transaction)
	if err != nil {
		return err
	}

	return nil
}
