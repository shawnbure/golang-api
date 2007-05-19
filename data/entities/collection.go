package entities

import "gorm.io/datatypes"

type Collection struct {
	ID                       uint64         `gorm:"primaryKey" json:"id"`
	Name                     string         `json:"name"`
	TokenID                  string         `json:"tokenId"`
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

	CreatorID uint64 `json:"creatorId"`
}
