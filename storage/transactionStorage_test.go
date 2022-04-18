package storage

import (
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
		lastFetchedId := int64(-1)
		lastTimestamp := int64(-1)
		howMuchRow := 2
		transactions, err := GetAllTransactionsWithPagination(lastFetchedId, lastTimestamp, howMuchRow)
		require.Nil(t, err)

		require.Equal(t, len(transactions), 2, "The returned transactions array length does not matched")
	})

	t.Run("Get Transactions With Detail with pagination", func(t *testing.T) {
		lastFetchedId := int64(107)
		lastTimestamp := int64(1586480400)
		howMuchRow := 2
		transactions, err := GetAllTransactionsWithPagination(lastFetchedId, lastTimestamp, howMuchRow)
		require.Nil(t, err)

		require.Equal(t, len(transactions), 2, "The returned transactions array length does not matched")
	})

}

func defaultTransaction() entities.Transaction {
	return entities.Transaction{
		Hash:         "hash",
		Type:         "test",
		PriceNominal: 1_000_000_000_000_000_000_000,
		SellerID:     1,
		BuyerID:      3,
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
	transaction.Timestamp = uint64(time.Date(2020, 04, 9, 1, 0, 0, 0, time.UTC).Unix())
	err := AddTransaction(&transaction)
	if err != nil {
		return err
	}

	transaction = defaultTransaction()
	transaction.Type = entities.BuyToken
	transaction.Hash = "my_unique_hash5"
	transaction.Timestamp = uint64(time.Date(2020, 04, 10, 1, 0, 0, 0, time.UTC).Unix())
	err = AddTransaction(&transaction)
	if err != nil {
		return err
	}

	transaction = defaultTransaction()
	transaction.Type = entities.WithdrawToken
	transaction.Hash = "my_unique_hash2"
	transaction.Timestamp = uint64(time.Date(2020, 04, 10, 1, 0, 0, 0, time.UTC).Unix())
	err = AddTransaction(&transaction)
	if err != nil {
		return err
	}

	transaction = defaultTransaction()
	transaction.Type = entities.AuctionToken
	transaction.Hash = "my_unique_hash3"
	transaction.Timestamp = uint64(time.Date(2020, 04, 9, 1, 0, 0, 0, time.UTC).Unix())
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
