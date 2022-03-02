package dtos

import (
	"gorm.io/datatypes"

	"github.com/ENFT-DAO/youbei-api/data/entities"
)

type CollectionToCheck struct {
	CollectionAddr string `json:"collectionAddr"`
	TokenID        string `json:"tokenId"`
	Counter        int    `json:"counter"`
}
type CollectionStatistics struct {
	ItemsTotal   uint64          `json:"itemsTotal"`
	OwnersTotal  uint64          `json:"ownersTotal"`
	FloorPrice   float64         `json:"floorPrice"`
	VolumeTraded float64         `json:"volumeTraded"`
	AttrStats    []AttributeStat `json:"attributes"`
}

type Attribute struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
}

type AttributeStat struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
	Total     uint64      `json:"total"`
}

type MetadataLinkResponse struct {
	Name       string      `json:"name"`
	Image      string      `json:"image"`
	Attributes []Attribute `json:"attributes"`
}

type ExtendedCollectionDto struct {
	Collection entities.Collection  `json:"collection"`
	Statistics CollectionStatistics `json:"statistics"`

	CreatorName          string `json:"creatorName"`
	CreatorWalletAddress string `json:"creatorWalletAddress"`
}

type CollectionCacheInfo struct {
	CollectionId    uint64         `json:"collectionId"`
	CollectionName  string         `json:"collectionName"`
	CollectionFlags datatypes.JSON `json:"collectionFlags"`
}
