package data

type Collection struct {
	ID            uint64 `gorm:"primaryKey"`
	Name          string
	TokenID       string
	Description   string
	Website       string
	DiscordLink   string
	TwitterLink   string
	InstagramLink string
	TelegramLink  string

	CreatorID uint64
}
