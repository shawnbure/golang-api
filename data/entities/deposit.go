package entities

type Deposit struct {
	ID            uint64  `gorm:"primaryKey" json:"id"`
	AmountNominal float64 `json:"amountNominal"`
	AmountString  string  `json:"amountString"`

	OwnerId uint64 `json:"ownerId"`
}
