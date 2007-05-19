package images

type AccountCoverImage struct {
	ID          uint64 `gorm:"primaryKey"`
	ImageBase64 string

	AccountID uint64
}
