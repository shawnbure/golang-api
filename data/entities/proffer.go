package entities

type Proffer struct {
	ID            uint64      `gorm:"primaryKey" json:"id"`
	Type          ProfferType `json:"type"`
	AmountNominal uint64      `json:"amountNominal"`
	AmountString  string      `json:"amountString"`
	Timestamp     uint64      `json:"timestamp"`
	TxHash        string      `json:"txHash"`

	AssetID   uint64 `json:"assetId"`
	OfferorID uint64 `json:"offerorId"`
}

type ProfferType string
