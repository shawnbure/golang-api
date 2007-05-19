package entities

type Offer struct {
	ID             uint64  `gorm:"primaryKey" json:"id"`
	AmountNominal  float64 `json:"amountNominal"`
	AmountString   string  `json:"amountString"`
	Expire         uint64  `json:"expire"`
	Timestamp      uint64  `json:"timestamp"`
	OfferorAddress string  `json:"offerorAddress"`
	TxHash         string  `json:"txHash"`

	TokenID uint64 `json:"tokenId"`
}
