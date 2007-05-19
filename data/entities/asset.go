package entities

import "gorm.io/datatypes"

type Asset struct {
	ID               uint64 `gorm:"primaryKey"`
	TokenID          string
	Nonce            uint64
	PriceNominal     float64
	RoyaltiesPercent float64
	MetadataLink     string
	CreatedAt        uint64
	Listed           bool
	Attributes       datatypes.JSON
	TokenName        string
	ImageLink        string
	Hash             string

	OwnerId      uint64
	CollectionID uint64
}
