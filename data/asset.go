package data

type Asset struct {
	ID               uint64 `gorm:"primaryKey"`
	TokenID          string
	Nonce            uint64
	Price            string
	RoyaltiesPercent uint64
	Link             string
	CreatedAt        uint64
	Listed           bool

	OwnerId      uint64
	CollectionID uint64
}
