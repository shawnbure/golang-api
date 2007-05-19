package images

type CollectionCoverImage struct {
	ID          uint64 `gorm:"primaryKey"`
	ImageBase64 string

	CollectionID uint64
}
