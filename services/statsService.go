package services

import (
	"fmt"
	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
)

func GetAllTransactionsWithPagination(args GetAllTransactionsWithPaginationArgs) ([]entities.TransactionDetail, int64, error) {
	total, err := storage.GetTransactionsCountWithCriteria(args.Filter, &entities.QueryFilter{})
	if err != nil {
		return nil, 0, err
	}

	transactions, err := storage.GetAllTransactionsWithPagination(args.LastTimestamp, args.CurrentPage, args.NextPage, args.Limit, args.Filter)
	if err != nil {
		return nil, 0, err
	}

	// Let's check the cache first
	localCacher := cache.GetLocalCacher()

	for index, item := range transactions {
		if item.TxType == string(entities.BuyToken) {
			// Get the buyer address from cache
			address, err := localCacher.Get(fmt.Sprintf(AddressByIdKeyFormat, item.ToId))
			if err == nil {
				item.ToAddress = address.(string)
				transactions[index] = item
			} else {
				// Get address from database
				acc, err := storage.GetAccountById(uint64(item.ToId))
				if err == nil {
					item.ToAddress = acc.Address
					transactions[index] = item

					// set the cache
					_ = localCacher.SetWithTTL(fmt.Sprintf(AddressByIdKeyFormat, item.ToId), acc.Address, AddressByIdExpirePeriod)
				}
			}
		}
	}

	if args.NextPage < args.CurrentPage {
		// reversing array
		for i, j := 0, len(transactions)-1; i < j; i, j = i+1, j-1 {
			transactions[i], transactions[j] = transactions[j], transactions[i]
		}
	}

	return transactions, total, nil
}
