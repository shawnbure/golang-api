package entities

type Transaction struct {
	ID           uint64 `gorm:"primaryKey"`
	Hash         string
	Type         TxType
	PriceNominal float64
	Timestamp    uint64

	SellerID     uint64
	BuyerID      uint64
	TokenID      uint64
	CollectionID uint64
}

type TxType string

const (
	ListToken     TxType = "List"
	BuyToken             = "Buy"
	WithdrawToken        = "Withdraw"
)
