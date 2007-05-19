package data

type Collection struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string
	TokenID     string
	Description string

	CreatorID   uint64
}
