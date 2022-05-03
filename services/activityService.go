package services

import (
	"fmt"
	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	"time"
)

const (
	AddressByIdKeyFormat    = "Address:ByDbId:%d"
	AddressByIdExpirePeriod = 24 * 2 * time.Hour
)

func GetAllActivities(args GetAllActivityArgs) ([]entities.Activity, error) {
	transactions, err := storage.GetAllActivitiesWithPagination(args.LastTimestamp, args.Limit, args.Filter)
	if err != nil {
		return nil, err
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

	return transactions, nil
}
