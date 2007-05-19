package images

type CollectionProfileImage struct {
	ID          uint64 `gorm:"primaryKey"`
	ImageBase64 string

	CollectionID uint64
}
