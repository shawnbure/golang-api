package data

type Asset struct {
	ID               uint64 `gorm:"primaryKey"`
	TokenID          string
	Nonce            uint64
	PriceNominal     float64
	RoyaltiesPercent float64
	Link             string
	CreatedAt        uint64
	Listed           bool
	LinkResponse	 string

	OwnerId      uint64
	CollectionID uint64
}
