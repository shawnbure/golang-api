package entities

type Transaction struct {
	ID           uint64  `gorm:"primaryKey" json:"id"`
	Hash         string  `json:"hash" gorm:"uniqueIndex"`
	Type         TxType  `json:"type" `
	PriceNominal float64 `json:"priceNominal"`
	Timestamp    uint64  `json:"timestamp"`
	SellerID     uint64  `json:"sellerId"`
	BuyerID      uint64  `json:"buyerId"`
	TokenID      uint64  `json:"tokenId"`
	CollectionID uint64  `json:"collectionId"`
}

type TxType string

const (
	ListToken     TxType = "List"
	BuyToken             = "Buy"
	WithdrawToken        = "Withdraw"
	AuctionToken         = "Auction"
)
