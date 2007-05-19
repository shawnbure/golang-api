package entities

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
	CreatedAt     uint64
	Priority      uint64

	CreatorID uint64
}
