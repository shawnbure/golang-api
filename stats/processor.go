package stats

import (
	"encoding/json"

	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/storage"
)

type CollectionMetadata struct {
	NumItems  uint64
	Owners    map[uint64]bool
	AttrStats map[string]map[string]int
}

func ComputeStatisticsForCollection(collectionId uint64) (*dtos.CollectionStatistics, error) {
	var stats dtos.CollectionStatistics

	//TODO: refactor this to something smarter. Min price is not good
	minPrice, err := storage.GetMinBuyPriceForTransactionsWithCollectionId(collectionId)
	if err != nil {
		return nil, err
	}

	sumPrice, err := storage.GetSumBuyPriceForTransactionsWithCollectionId(collectionId)
	if err != nil {
		return nil, err
	}

	collectionMetadata, err := ComputeCollectionMetadata(collectionId)
	if err != nil {
		return nil, err
	}

	stats = dtos.CollectionStatistics{
		ItemsTotal:   collectionMetadata.NumItems,
		OwnersTotal:  uint64(len(collectionMetadata.Owners)),
		FloorPrice:   minPrice,
		VolumeTraded: sumPrice,
	}

	return &stats, nil
}

func ComputeCollectionMetadata(collectionId uint64) (*CollectionMetadata, error) {
	offset := 0
	limit := 1_000
	numItems := 0
	ownersIDs := make(map[uint64]bool)
	attrStats := make(map[string]map[string]int)

	for {
		tokens, innerErr := storage.GetListedTokensByCollectionIdWithOffsetLimit(collectionId, offset, limit)
		if innerErr != nil {
			return nil, innerErr
		}
		if len(tokens) == 0 {
			break
		}

		numItems = numItems + len(tokens)
		for _, token := range tokens {
			tokenAttrs := make(map[string]string)
			ownersIDs[token.OwnerId] = true

			innerErr = json.Unmarshal(token.Attributes, &tokenAttrs)
			if innerErr != nil {
				continue
			}

			for attrName, attrValue := range tokenAttrs {
				if _, ok := attrStats[attrName]; ok {
					attrStats[attrName][attrValue] += 1
				} else {
					attrStats[attrName] = map[string]int{attrValue: 1}
				}
			}
		}

		offset = limit
		limit = limit + 1_000
	}

	result := CollectionMetadata{
		NumItems:  uint64(numItems),
		Owners:    ownersIDs,
		AttrStats: attrStats,
	}
	return &result, nil
}
