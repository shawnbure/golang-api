package entities

type CollectionIndexer struct {
	ID             uint64 `json:"id" gorm:"primaryKey"`
	CollectionAddr string `json:"collectionAddr"`
	CollectionName string `json:"collectionName"`
	LastIndex      uint64 `json:"lastIndex"`
	LastNonce      uint64 `json:"lastNonce"`
	UpdatedAt      int64  `json:"updatedAt" gorm:"autoUpdateTime:milli"`  // Set to current unix seconds on updaing or if it is zero on creating
	CreatedAt      int64  `json:"created_at" gorm:"autoCreateTime:milli"` // Use unix seconds as creating time
}
