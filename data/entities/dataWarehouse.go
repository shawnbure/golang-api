package entities

type AggregatedVolumePerHour struct {
	ID             uint64  `gorm:"primaryKey" json:"id"`
	Hour           int64   `json:"hour" gorm:"index:,unique"`
	BuyVolume      float64 `json:"buyVolume"`
	ListVolume     float64 `json:"listVolume"`
	WithdrawVolume float64 `json:"withdrawVolume"`
}

type AggregatedVolumePerCollectionPerHour struct {
	ID           uint64 `gorm:"primaryKey" json:"id"`
	Hour         int64  `json:"hour" gorm:"uniqueIndex:uidx_aggregated_collection_volume_per_hour"`
	CollectionId uint64 `json:"collectionId" gorm:"uniqueIndex:uidx_aggregated_collection_volume_per_hour"`
	//Collection     Collection `json:"collection"`
	BuyVolume      float64 `json:"buyVolume"`
	ListVolume     float64 `json:"listVolume"`
	WithdrawVolume float64 `json:"withdrawVolume"`
}

type GroupAggregatedVolumePerCollection struct {
	Total        float64 `json:"total"`
	CollectionId uint64  `json:"collectionId"`
	Type         string  `json:"type"`
}
