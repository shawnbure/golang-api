package dtos

import (
	"encoding/json"

	"github.com/erdsea/erdsea-api/data/entities"
)

type ExtendedTokenDto struct {
	entities.Token

	AttributesMap  map[string]string
	CollectionID   string
	CollectionName string

	OwnerName          string
	OwnerWalletAddress string

	Stats CollectionStatistics
}

func CreateExtendedTokenDto(
	token entities.Token,
	collectionID string,
	collectionName string,
	ownerName string,
	ownerWalletAddress string,
	collStats CollectionStatistics,
) (*ExtendedTokenDto, error) {
	e := &ExtendedTokenDto{
		Token:              token,
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
