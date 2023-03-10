package stats

import (
	"encoding/json"
	"fmt"

	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/storage"
)

type CollectionMetadata struct {
	NumItems  uint64
	Owners    map[uint64]bool
	AttrStats []dtos.AttributeStat
}

func ComputeStatisticsForCollection(collectionId uint64) (*dtos.CollectionStatistics, error) {
	var stats dtos.CollectionStatistics

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
		AttrStats:    collectionMetadata.AttrStats,
	}

	return &stats, nil
}

func ComputeCollectionMetadata(collectionId uint64) (*CollectionMetadata, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	offset := 0
	limit := 1_000
	numItems := 0
	ownersIDs := make(map[uint64]bool)
	var globalAttrs []dtos.AttributeStat
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
			tokenAttrs := []map[string]interface{}{}
			ownersIDs[token.OwnerID] = true
			innerErr = json.Unmarshal(token.Attributes, &tokenAttrs)
			if innerErr != nil {
				continue
			}
			if len(tokenAttrs) == 0 {
				continue
			}
			for _, obj := range tokenAttrs {
				attributeFound := false
				for index, globalAttr := range globalAttrs {

					if _, ok := obj["trait_type"]; !ok {
						for k, v := range obj {
							if globalAttr.TraitType == k && globalAttr.Value == v {
								attributeFound = true
								globalAttrs[index].Total++
							}
						}
					} else {
						if globalAttr.TraitType == obj["trait_type"] && globalAttr.Value == obj["value"] {
							attributeFound = true
							globalAttrs[index].Total++
						}
					}
				}

				if !attributeFound {
					if _, ok := obj["trait_type"]; !ok {
						for k, v := range obj {
							globalAttrs = append(globalAttrs, dtos.AttributeStat{
								TraitType: k,
								Value:     v,
								Total:     1,
							})
						}
					} else {
						globalAttrs = append(globalAttrs, dtos.AttributeStat{
							TraitType: obj["trait_type"].(string),
							Value:     obj["value"],
							Total:     1,
						})

					}
				}
			}
		}

		offset += limit
	}

	result := CollectionMetadata{
		NumItems:  uint64(numItems),
		Owners:    ownersIDs,
		AttrStats: globalAttrs,
	}
	return &result, nil
}
