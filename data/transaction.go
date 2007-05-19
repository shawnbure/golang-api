package data

type Transaction struct {
	ID           uint64 `gorm:"primaryKey"`
	Hash         string
	Type         TxType
	PriceNominal float64
	Timestamp    uint64

	SellerID uint64
	BuyerID  uint64
	AssetID  uint64
}

type TxType string

const (
	ListAsset     TxType = "List"
	BuyAsset             = "Buy"
	WithdrawAsset        = "Withdraw"
)
