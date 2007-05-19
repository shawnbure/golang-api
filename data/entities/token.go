package entities

import "gorm.io/datatypes"

type Token struct {
	ID               uint64         `gorm:"primaryKey" json:"id"`
	TokenID          string         `json:"tokenId"`
	Nonce            uint64         `json:"nonce"`
	PriceString      string         `json:"priceString"`
	PriceNominal     float64        `json:"priceNominal"`
	RoyaltiesPercent float64        `json:"royaltiesPercent"`
	MetadataLink     string         `json:"metadataLink"`
	CreatedAt        uint64         `json:"createdAt"`
	Listed           bool           `json:"listed"`
	Attributes       datatypes.JSON `json:"attributes"`
	TokenName        string         `json:"tokenName"`
	ImageLink        string         `json:"imageLink"`
	Hash             string         `json:"hash"`

	OwnerId      uint64 `json:"ownerId"`
	CollectionID uint64 `json:"collectionId"`
}
