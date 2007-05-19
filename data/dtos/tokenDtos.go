package dtos

import (
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
	entities.Token      `json:"token"`
	entities.Collection `json:"collection"`

	OwnerName          string `json:"ownerName"`
	OwnerWalletAddress string `json:"ownerWalletAddress"`

	CollectionStats CollectionStatistics `json:"collectionStats"`
}

func CreateExtendedTokenDto(
	token entities.Token,
	collection entities.Collection,
	ownerName string,
	ownerWalletAddress string,
	collStats CollectionStatistics,
) (*ExtendedTokenDto, error) {
	e := &ExtendedTokenDto{
		Token:              token,
		Collection:         collection,
		OwnerName:          ownerName,
		OwnerWalletAddress: ownerWalletAddress,
		CollectionStats:    collStats,
	}

	return e, nil
}
