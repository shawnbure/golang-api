package entities

import "gorm.io/datatypes"

type Token struct {
	ID                  uint64         `gorm:"primaryKey" json:"id"`
	TokenID             string         `json:"tokenId"`
	Nonce               uint64         `json:"nonce"`
	NonceStr            string         `json:"nonceStr"`
	PriceString         string         `json:"priceString"`
	PriceNominal        float64        `json:"priceNominal"`
	RoyaltiesPercent    float64        `json:"royaltiesPercent"`
	MetadataLink        string         `json:"metadataLink"`
	CreatedAt           uint64         `json:"createdAt"`
	Status              TokenStatus    `json:"state"`
	Attributes          datatypes.JSON `json:"attributes"`
	TokenName           string         `json:"tokenName"`
	ImageLink           string         `json:"imageLink"`
	Hash                string         `json:"hash"`
	MintTxHash          string         `json:"mintTxHash"`
	LastBuyPriceNominal float64        `json:"lastBuyPriceNominal"`
	AuctionStartTime    uint64         `json:"auctionStartTime"`
	AuctionDeadline     uint64         `json:"auctionDeadline"`
	OnSale              bool           `json:"onSale"`
	OwnerId             uint64         `json:"ownerId"`
	CollectionID        uint64         `json:"collectionId"`
}

type TokenStatus string

const (
	List    TokenStatus = "List"
	Auction             = "Auction"
	None                = "None"
)
