package dtos

type TradesCount struct {
	Total    int64 `json:"Total"`
	Buy      int64 `json:"Buy"`
	Withdraw int64 `json:"Withdraw"`
}

type TradesVolumeTotal struct {
	Sum        string `json:"sum"`
	LastUpdate int64  `json:"last_update"`
}

type TradesVolume struct {
	Sum  string `json:"sum"`
	Date string `json:"date"`
}
