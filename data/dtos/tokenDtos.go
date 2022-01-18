package dtos

import (
	"github.com/ENFT-DAO/youbei-api/data/entities"
)

type ExtendedTokenDto struct {
	entities.Token     `json:"token"`
	OwnerName          string `json:"ownerName"`
	OwnerWalletAddress string `json:"ownerWalletAddress"`
}

type OwnedTokenDto struct {
	entities.Token      `json:"token"`
	CollectionCacheInfo `json:"collection"`
}
