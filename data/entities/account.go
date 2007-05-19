package entities

type Account struct {
	ID            uint64 `gorm:"primaryKey"`
	Address       string
	Name          string
	Description   string
	Website       string
	TwitterLink   string
	InstagramLink string
	CreatedAt     uint64
}
