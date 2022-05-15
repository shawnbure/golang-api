package entities

type AggregatedVolumePerHour struct {
	ID             uint64  `gorm:"primaryKey" json:"id"`
	Hour           int64   `json:"hour" gorm:"index:,unique"`
	BuyVolume      float64 `json:"buyVolume"`
	ListVolume     float64 `json:"listVolume"`
	WithdrawVolume float64 `json:"withdrawVolume"`
}
