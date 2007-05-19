package entities

type Bid struct {
	ID               uint64  `gorm:"primaryKey" json:"id"`
	BidAmountNominal float64 `json:"bidAmountNominal"`
	BidAmountString  string  `json:"bidAmountString"`
	Timestamp        uint64  `json:"timestamp"`
	BidderAddress    string  `json:"bidderAddress"`
	TxHash           string  `json:"txHash"`

	TokenID uint64 `json:"tokenId"`
}
