package entities

type TopVolumeByAddress struct {
	FromTime string  `json:"from_time"`
	ToTime   string  `json:"to_time"`
	Address  string  `json:"address"`
	Volume   float64 `json:"volume"`
}
