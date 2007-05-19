package data

type Asset struct {
	ID      uint64 `gorm:"primaryKey"`
	TokenID string
	Nonce   uint64
	Price   string
	Link    string
	Listed  bool

	OwnerId      uint64
	CollectionID uint64
}
