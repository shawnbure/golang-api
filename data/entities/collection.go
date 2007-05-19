package entities

type Collection struct {
	ID            uint64 `gorm:"primaryKey"`
	Name          string `json:"name"`
	TokenID       string `json:"tokenId"`
	Description   string `json:"description"`
	Website       string `json:"website"`
	DiscordLink   string `json:"discordLink"`
	TwitterLink   string `json:"twitterLink"`
	InstagramLink string `json:"instagramLink"`
	TelegramLink  string `json:"telegramLink"`
	CreatedAt     uint64 `json:"createdAt"`
	Priority      uint64 `json:"priority"`

	CreatorID uint64 `json:"creatorId"`
}
