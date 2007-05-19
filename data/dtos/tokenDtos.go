package dtos

import (
	"github.com/erdsea/erdsea-api/data/entities"
)

type ExtendedTokenDto struct {
	entities.Token      `json:"token"`
	entities.Collection `json:"collection"`

	OwnerName          string `json:"ownerName"`
	OwnerWalletAddress string `json:"ownerWalletAddress"`

	CreatorName          string `json:"creatorName"`
	CreatorWalletAddress string `json:"creatorWalletAddress"`

	CollectionStats CollectionStatistics `json:"collectionStats"`
}

func CreateExtendedTokenDto(
	token entities.Token,
	collection entities.Collection,
	ownerName string,
	ownerWalletAddress string,
	creatorName string,
	creatorWalletAddress string,
	collStats CollectionStatistics,
) (*ExtendedTokenDto, error) {
	e := &ExtendedTokenDto{
		Token:                token,
		Collection:           collection,
		OwnerName:            ownerName,
		OwnerWalletAddress:   ownerWalletAddress,
		CreatorName:          creatorName,
		CreatorWalletAddress: creatorWalletAddress,
		CollectionStats:      collStats,
	}

	return e, nil
}
