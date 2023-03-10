package entities

import "gorm.io/datatypes"

type Token struct {
	ID                   uint64         `gorm:"primaryKey" json:"id"`
	TokenID              string         `json:"tokenId" gorm:"index:token_nonces,unique;not null"`
	Nonce                uint64         `json:"nonce" gorm:"index:token_nonces,unique;not null"`
	NonceStr             string         `json:"nonceStr" gorm:"index:token_nonces,unique;not null"`
	PriceString          string         `json:"priceString"`
	PriceNominal         float64        `json:"priceNominal"`
	RoyaltiesPercent     float64        `json:"royaltiesPercent"`
	MetadataLink         string         `json:"metadataLink"`
	CreatedAt            uint64         `json:"createdAt"`
	UpdatedAt            int64          `json:"updatedAt" gorm:"autoUpdateTime:milli"` // Set to current unix seconds on updaing or if it is zero on creating
	Status               TxType         `json:"state"`
	Attributes           datatypes.JSON `json:"attributes"`
	TokenName            string         `json:"tokenName"`
	ImageLink            string         `json:"imageLink"`
	Hash                 string         `json:"hash"`
	MintTxHash           string         `json:"mintTxHash"`
	LastBuyPriceNominal  float64        `json:"lastBuyPriceNominal"`
	AuctionStartTime     uint64         `json:"auctionStartTime"`
	AuctionDeadline      uint64         `json:"auctionDeadline"`
	OnSale               bool           `json:"onSale"`
	OwnerID              uint64         `json:"ownerId"`
	CollectionID         uint64         `json:"collectionId"`
	Collection           Collection     `json:"collection"`
	LastMarketTimestamp  uint64         `json:"lastMarketTimestamp"`
	Owner                Account        `json:"owner"`
	TxConfirmed          bool           `json:"txConfirmed"`
	OnStake              bool           `json:"onStake" gorm:"default:false"`
	StakeDate            uint64         `json:"stakeDate"`
	StakeType            string         `json:"stakeType"`
	RarityUsedTraitCount uint           `json:"rarityUsedTraitCount" gorm:"default:0"`
	Rank                 uint           `json:"rank" gorm:"default:0"`
	RarityScore          float64        `json:"rarityScore" gorm:"default:0.0"`
	RarityScoreNorm      float64        `json:"rarityScoreNorm" gorm:"default:0.0"`
	IsRarityInserted     bool           `json:"isRarityInserted" gorm:"default:false"`
	RarityLastUpdated    uint64         `json:"rarityLastUpdated" gorm:"autoUpdateTime:milli;default:0"`
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

// type TokenStatus string

// const (
// 	List    TokenStatus = "List"
// 	Stake   TokenStatus = "Stake"
// 	Auction TokenStatus = "Auction"
// 	None    TokenStatus = "None"
// )

type StakeType string

const (
	UBI StakeType = "UBI"
	DAO StakeType = "DAO"
)

type TokenExplorer struct {
	Token Token `json:"token" gorm:"embedded"`
	//TokenId                  string  `json:"tokenId"`
	//TokenStatus              string  `json:"tokenStatus"`
	//TokenName                string  `json:"tokenName"`
	//TokenImage               string  `json:"tokenImage"`
	//TokenAuctionStartTime    int64   `json:"tokenAuctionStartTime"`
	//TokenAuctionDeadline     int64   `json:"tokenAuctionDeadline"`
	//TokenCreatedAt           int64   `json:"tokenCreatedAt"`
	//TokenLastMarketTimestamp int64   `json:"tokenLastMarketTimestamp"`
	//TokenLastBuyPriceNominal float64 `json:"tokenLastBuyPriceNominal"`
	//TokenPriceNominal        float64 `json:"tokenPriceNominal"`
	Owner        Account    `json:"owner"`
	OwnerID      uint64     `json:"ownerId"`
	Collection   Collection `json:"collection"`
	CollectionID uint64     `json:"collectionId"`
}
