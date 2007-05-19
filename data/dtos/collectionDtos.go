package dtos

import "github.com/erdsea/erdsea-api/data/entities"

type CollectionStatistics struct {
	ItemsTotal   uint64          `json:"itemsTotal"`
	OwnersTotal  uint64          `json:"ownersTotal"`
	FloorPrice   float64         `json:"floorPrice"`
	VolumeTraded float64         `json:"volumeTraded"`
	AttrStats    []AttributeStat `json:"attributes"`
}

type Attribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

type AttributeStat struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
	Total     uint64 `json:"total"`
}

type ExtendedCollectionDto struct {
	Collection entities.Collection  `json:"collection"`
	Statistics CollectionStatistics `json:"statistics"`

	CreatorName          string `json:"creatorName"`
	CreatorWalletAddress string `json:"creatorWalletAddress"`
}
