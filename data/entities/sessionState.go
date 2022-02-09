package entities

type SessionState struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	AccountID uint64 `json:"accountId"`
	StateType uint64 `json:"stateType"`
	JsonData  string `json:"jsonData"`
	CreatedAt uint64 `json:"createdAt"`
}

const (
	SessionState_type_none              = 0
	SessionState_type_create_collection = 1
)
