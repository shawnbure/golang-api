package entities

type Account struct {
	ID            uint64 `gorm:"primaryKey" json:"id"`
	Address       string `json:"address"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Website       string `json:"website"`
	TwitterLink   string `json:"twitterLink"`
	InstagramLink string `json:"instagramLink"`
	CreatedAt     uint64 `json:"createdAt"`
}
