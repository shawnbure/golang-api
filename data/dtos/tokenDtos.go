package dtos

import (
	"github.com/erdsea/erdsea-api/data/entities"
)

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
