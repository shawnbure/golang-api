package dtos

import "github.com/erdsea/erdsea-api/data/entities"

type CollectionStatistics struct {
	ItemsTotal   uint64                    `json:"itemsTotal"`
	OwnersTotal  uint64                    `json:"ownersTotal"`
	FloorPrice   float64                   `json:"floorPrice"`
	VolumeTraded float64                   `json:"volumeTraded"`
	AttrStats    map[string]map[string]int `json:"attributes"`
}

type ExtendedCollectionDto struct {
	Collection entities.Collection  `json:"collection"`
	Statistics CollectionStatistics `json:"statistics"`
}
