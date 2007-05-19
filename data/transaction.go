package data

type Transaction struct {
	ID    uint64 `gorm:"primaryKey"`
	Hash  string
	Type  string
	Price string

	SellerID uint64
	BuyerID  uint64
	AssetID  uint64
}
