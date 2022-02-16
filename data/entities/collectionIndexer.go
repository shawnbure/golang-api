package entities

type CollectionIndexer struct {
	ID             uint64 `json:"id" gorm:"primaryKey"`
	CollectionAddr string `json:"collectionAddr"`
	LastIndex      uint64 `json:"lastIndex"`
}
