package dtos

import "github.com/ENFT-DAO/youbei-api/data/entities"

type BidDto struct {
	entities.Bid `json:"bid"`
	BidderName   string `json:"bidderName"`
}
