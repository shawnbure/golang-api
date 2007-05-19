package entities

type Proffer struct {
	ID            uint64      `gorm:"primaryKey" json:"id"`
	Type          ProfferType `json:"type"`
	AmountNominal float64     `json:"amountNominal"`
	AmountString  string      `json:"amountString"`
	Expire        uint64      `json:"expire"`
	Timestamp     uint64      `json:"timestamp"`
	TxHash        string      `json:"txHash"`

	TokenID   uint64 `json:"tokenId"`
	OfferorID uint64 `json:"offerorId"`
}

type ProfferType string

const (
	Bid   ProfferType = "Bid"
	Offer             = "Offer"
)
