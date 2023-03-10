package entities

type Account struct {
	ID               uint64      `gorm:"primaryKey" json:"id"`
	Address          string      `json:"address" gorm:"index:,unique"`
	Name             string      `json:"name" gorm:"default:random()::text;index:,unique"`
	Description      string      `json:"description"`
	Website          string      `json:"website"`
	TwitterLink      string      `json:"twitterLink"`
	InstagramLink    string      `json:"instagramLink"`
	CreatedAt        uint64      `json:"createdAt"`
	ProfileImageLink string      `json:"profileImageLink"`
	CoverImageLink   string      `json:"coverImageLink"`
	MintedCount      uint64      `json:"MintedCount" gorm:"default:0"`
	MaxBatchMint     uint64      `json:"maxBatchMint" gorm:"default:10"`
	MaxLifetimeMint  uint64      `json:"maxLifetimeMint" gorm:"default:10000"`
	Role             AccountRole `json:"role" gorm:"default:'RoleUser'"`
}

type AccountRole string

const (
	RoleUser  AccountRole = "RoleUser"
	RoleAdmin AccountRole = "RoleAdmin"
)
