package services

import (
	"fmt"
	"time"

	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
)

const (
	AddressByIdKeyFormat    = "Address:ByDbId:%d"
	AddressByIdExpirePeriod = 24 * 2 * time.Hour
)

func GetAllActivities(args GetAllActivityArgs) ([]entities.Activity, int64, error) {
	total, err := storage.GetTransactionsCountWithCriteria(args.CollectionFilter)
	if err != nil {
		return nil, 0, err
	}

	transactions, err := storage.GetAllActivitiesWithPagination(args.LastTimestamp, args.CurrentPage, args.NextPage, args.Limit, args.Filter, args.CollectionFilter)
	if err != nil {
		return nil, 0, err
	}

	// Let's check the cache first
	localCacher := cache.GetLocalCacher()

	for index, item := range transactions {
		if string(item.Transaction.Type) == string(entities.BuyToken) {
			// Get the buyer address from cache
			// address, err := localCacher.Get(fmt.Sprintf(AddressByIdKeyFormat, item.ToId))
			if err == nil {
				transactions[index] = item
			} else {
				// Get address from database
				acc, err := storage.GetAccountById(uint64(item.Transaction.BuyerID))
				if err == nil {
					transactions[index] = item

					// set the cache
					_ = localCacher.SetWithTTL(fmt.Sprintf(AddressByIdKeyFormat, item.Transaction.BuyerID), acc.Address, AddressByIdExpirePeriod)
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
