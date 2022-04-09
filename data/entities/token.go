package entities

import "gorm.io/datatypes"

type Token struct {
	ID                  uint64         `gorm:"primaryKey" json:"id"`
	TokenID             string         `json:"tokenId" gorm:"uniqueIndex:token_nonces;not null"`
	Nonce               uint64         `json:"nonce" gorm:"uniqueIndex:token_nonces;not null"`
	NonceStr            string         `json:"nonceStr" gorm:"uniqueIndex:token_nonces;not null"`
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
	LastMarketTimestamp uint64         `json:"lastMarketTimestamp"`
	Owner               Account        `json:"owner"`
	TxConfirmed         bool           `json:"txConfirmed"`
}

type TokenBC struct {
	Identifier           string      `json:"identifier"`
	Collection           string      `json:"collection"`
	Timestamp            uint64      `json:"timestamp"`
	Attributes           string      `json:"attributes"`
	Nonce                uint64      `json:"nonce"`
	Type                 string      `json:"type"`
	Name                 string      `json:"name"`
	Creator              string      `json:"creator"`
	Owner                string      `json:"owner"`
	Royalties            interface{} `json:"royalties"`
	URIs                 []string    `json:"uris"`
	URL                  string      `json:"url"`
	Media                interface{} `json:"media"`
	IsWhitelistedStorage bool        `json:"isWhitelistedStorage"`
	Metadata             interface{} `json:"metadata"`
	/*
		  "thumbnailUrl": "string",
		  "tags": [
		    "string"
		  ],
		  "metadata": {},
		  "owner": {},
		  "balance": {},
		  "supply": {},
		  "decimals": {},
		  "assets": {
		    "website": "string",
		    "description": "string",
		    "status": "string",
		    "pngUrl": "string",
		    "svgUrl": "string",
		    "lockedAccounts": {}
		  },
		  "ticker": "string",
		  "scamInfo": {}
		}
	*/
}

type TokenStatus string

const (
	List    TokenStatus = "List"
	Auction             = "Auction"
	None                = "None"
)
