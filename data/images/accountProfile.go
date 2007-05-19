package images

type AccountProfileImage struct {
	ID          uint64 `gorm:"primaryKey"`
	ImageBase64 string

	AccountID uint64
}
