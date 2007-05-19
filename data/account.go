package data

type Account struct {
	ID      uint64 `gorm:"primaryKey"`
	Address string
}
