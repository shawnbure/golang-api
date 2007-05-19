package data

type Asset struct {
	ID      uint64 `gorm:"primaryKey"`
	TokenID string
	Nonce   uint64
	Price   uint64
	Link    string

	CreatorID    uint64
	CollectionID uint64
}
