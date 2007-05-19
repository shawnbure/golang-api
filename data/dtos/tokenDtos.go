package dtos

import (
	"encoding/json"

	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/services"
)

type ExtendedTokenDto struct {
	entities.Asset

	AttributesMap  map[string]string
	CollectionID   string
	CollectionName string

	OwnerName          string
	OwnerWalletAddress string

	Stats services.CollectionStatistics
}

func CreateExtendedTokenDto(
	asset entities.Asset,
	collectionID string,
	collectionName string,
	ownerName string,
	ownerWalletAddress string,
	collStats services.CollectionStatistics,
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
