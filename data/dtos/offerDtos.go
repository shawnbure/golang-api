package dtos

import "github.com/ENFT-DAO/youbei-api/data/entities"

type OfferDto struct {
	entities.Offer `json:"offer"`
	OfferorName    string `json:"offerorName"`
}
