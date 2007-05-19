package entities

type Transaction struct {
	ID           uint64  `gorm:"primaryKey" json:"id"`
	Hash         string  `json:"hash"`
	Type         TxType  `json:"type"`
	PriceNominal float64 `json:"price_nominal"`
	Timestamp    uint64  `json:"timestamp"`

	SellerID     uint64 `json:"seller_id"`
	BuyerID      uint64 `json:"buyer_id"`
	TokenID      uint64 `json:"token_id"`
	CollectionID uint64 `json:"collection_id"`
}

type TxType string

const (
	ListToken     TxType = "List"
	BuyToken             = "Buy"
	WithdrawToken        = "Withdraw"
)
