package entities

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type Collection struct {
	ID                       uint64         `gorm:"primaryKey" json:"id"`
	Name                     string         `json:"name"`
	CollectionTokenID        string         `json:"collectionTokenId" gorm:"index:,unique"`
	Description              string         `json:"description"`
	Website                  string         `json:"website"`
	DiscordLink              string         `json:"discordLink"`
	TwitterLink              string         `json:"twitterLink"`
	InstagramLink            string         `json:"instagramLink"`
	TelegramLink             string         `json:"telegramLink"`
	CreatedAt                uint64         `json:"createdAt" gorm:"autoCreateTime:milli"`
	Priority                 uint64         `json:"priority"`
	ContractAddress          string         `json:"contractAddress"` // could be contract address of creator role
	MintPricePerTokenString  string         `json:"mintPricePerTokenString"`
	MintPricePerTokenNominal float64        `json:"mintPricePerTokenNominal"`
	Flags                    datatypes.JSON `json:"flags"`
	ProfileImageLink         string         `json:"profileImageLink"`
	CoverImageLink           string         `json:"coverImageLink"`
	IsVerified               bool           `json:"isVerified"`
	IsStakeable              bool           `json:"isStakeable" gorm:"default:false"`
	Type                     uint64         `json:"type"`
	TokenBaseURI             string         `json:"tokenBaseURI"`
	MaxSupply                uint64         `json:"maxSupply"`
	MetaDataBaseURI          string         `json:"metaDataBaseURI"`
	MintStartDate            uint64         `json:"mintStartDate"`
	MintEndDate              uint64         `json:"mintEndDate"`
	CreatorID                uint64         `json:"creatorId"`

	//`gorm:"type:bool;default:false"`

	//AccountName              string `json:"accountName"`
	//AcccountProfileImageLink string `json:"accountProfileImageLink"`
}

const (
	Collection_type_none        = 0
	Collection_type_whitelisted = 1
	Collection_type_noteworthy  = 2
)

type DeployerStat struct {
	LastIndex    uint64 `json:"lastIndex"`
	DeployerAddr string `json:"deployerAddr"`
	ID           uint64 `json:"id" gorm:"primaryKey"`
	UpdatedAt    int64  `json:"updatedAt" gorm:"autoUpdateTime:milli"`  // Set to current unix seconds on updaing or if it is zero on creating
	CreatedAt    int64  `json:"created_at" gorm:"autoCreateTime:milli"` // Use unix seconds as creating time
}
type CollectionBC struct {
	Identifier           string         `json:"identifier" `
	Collection           string         `json:"collection"`
	Nonce                uint64         `json:"nonce"`
	NFTType              string         `json:"type"`
	Creator              string         `json:"creator"`
	Royalties            uint64         `json:"royalties"`
	URIs                 pq.StringArray `json:"uris"`
	URL                  string         `json:"url"`
	IsWhitelistedStorage bool           `json:"isWhitelistedStorage"`
	Metadata             JSONB          `json:"metadata"`
	Ticker               string         `json:"ticker"`
}

type CollectionAccount struct {
	ID                       uint64         `gorm:"primaryKey" json:"id"`
	Name                     string         `json:"name"`
	TokenID                  string         `json:"tokenId" `
	Description              string         `json:"description"`
	Website                  string         `json:"website"`
	DiscordLink              string         `json:"discordLink"`
	TwitterLink              string         `json:"twitterLink"`
	InstagramLink            string         `json:"instagramLink"`
	TelegramLink             string         `json:"telegramLink"`
	CreatedAt                uint64         `json:"createdAt"`
	Priority                 uint64         `json:"priority"`
	ContractAddress          string         `json:"contractAddress"`
	MintPricePerTokenString  string         `json:"mintPricePerTokenString"`
	MintPricePerTokenNominal float64        `json:"mintPricePerTokenNominal"`
	Flags                    datatypes.JSON `json:"flags"`
	ProfileImageLink         string         `json:"profileImageLink"`
	CoverImageLink           string         `json:"coverImageLink"`
	IsVerified               bool           `json:"isVerified"`
	Type                     uint64         `json:"type"`

	CreatorID uint64 `json:"creatorId"`

	AccountName              string `json:"accountName"`
	AcccountProfileImageLink string `json:"accountProfileImageLink"`
	AccountAddress           string `json:"accountAddress"`
}
