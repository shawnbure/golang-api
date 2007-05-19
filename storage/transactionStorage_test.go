package storage

import (
	"testing"

	"github.com/erdsea/erdsea-api/data"
	"github.com/stretchr/testify/require"
)

func Test_AddTransaction(t *testing.T) {
	connectToTestDb()

	transaction := defaultTransaction()
	err := AddTransaction(&transaction)
	require.Nil(t, err)

	var transactionRead data.Transaction
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

func Test_GetTransactionsByAssetId(t *testing.T) {
	connectToTestDb()

	transaction := defaultTransaction()
	err := AddTransaction(&transaction)
	require.Nil(t, err)

	otherTransaction := defaultTransaction()
	err = AddTransaction(&otherTransaction)
	require.Nil(t, err)

	transactionsRead, err := GetTransactionsByAssetId(transaction.TokenID)
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

func defaultTransaction() data.Transaction {
	return data.Transaction{
		Hash:         "hash",
		Type:         "test",
		PriceNominal: 1_000_000_000_000_000_000_000,
		SellerID:     1,
		BuyerID:      2,
		TokenID:      3,
		CollectionID: 1,
	}
}
