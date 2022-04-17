package dtos

import "github.com/ENFT-DAO/youbei-api/data/entities"

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

type TokensTotalCount struct {
	Sum int64 `json:"sum"`
}

type StatTransactionsList struct {
	Transactions []entities.TransactionDetail `json:"transactions"`
}
