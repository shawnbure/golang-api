package dtos

import "github.com/erdsea/erdsea-api/data/entities"

type BidDto struct {
	entities.Bid `json:"bid"`
	BidderName   string `json:"bidderName"`
}
