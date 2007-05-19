package dtos

import (
	"encoding/json"

	"github.com/erdsea/erdsea-api/data/entities"
)

type CollectionStatistics struct {
	ItemsTotal   uint64                    `json:"itemsTotal"`
	OwnersTotal  uint64                    `json:"ownersTotal"`
	FloorPrice   float64                   `json:"floorPrice"`
	VolumeTraded float64                   `json:"volumeTraded"`
	AttrStats    map[string]map[string]int `json:"attributes"`
}

type ExtendedTokenDto struct {
	entities.Asset

	AttributesMap  map[string]string
	CollectionID   string
	CollectionName string

	OwnerName          string
	OwnerWalletAddress string

	Stats CollectionStatistics
}

func CreateExtendedTokenDto(
	asset entities.Asset,
	collectionID string,
	collectionName string,
	ownerName string,
	ownerWalletAddress string,
	collStats CollectionStatistics,
) (*ExtendedTokenDto, error) {
	e := &ExtendedTokenDto{
		Asset:              asset,
		CollectionID:       collectionID,
		CollectionName:     collectionName,
		OwnerName:          ownerName,
		OwnerWalletAddress: ownerWalletAddress,
		Stats:              collStats,
	}

	err := json.Unmarshal(e.Attributes, &e.AttributesMap)
	if err != nil {
		return nil, err
	}

	return e, nil
}
