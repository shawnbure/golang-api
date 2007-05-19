package data

type Asset struct {
	ID      uint64 `gorm:"primaryKey"`
	TokenID string
	Nonce   uint64
	Price   string //no big.Int support in gorm
	Link    string

	OwnerId      uint64
	CollectionID uint64
}
