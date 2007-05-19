package dtos

import "github.com/erdsea/erdsea-api/data/entities"

type OfferDto struct {
	entities.Offer `json:"offer"`
	OfferorName    string `json:"offerorName"`
}
