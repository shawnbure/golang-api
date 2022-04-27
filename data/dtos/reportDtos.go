package dtos

import "github.com/ENFT-DAO/youbei-api/data/entities"

type ReportLast24HoursTransactionsList struct {
	Transactions []entities.TransactionDetail `json:"transactions"`
}

type ReportLast24HoursOverall struct {
	FromTime          string  `json:"from_time"`
	ToTime            string  `json:"to_time"`
	TotalVolume       float64 `json:"total_volume"`
	TotalVolumeStr    string  `json:"total_volume_str"`
	TotalTransactions int     `json:"total_transactions"`
}

type ReportTopVolumeByAddress struct {
	Records []entities.TopVolumeByAddress `json:"records"`
}

type ReportTopVolumeByAddressTransactionsList struct {
	Transactions []entities.TransactionDetail `json:"transactions"`
}
