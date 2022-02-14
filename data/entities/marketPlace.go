package entities

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/lib/pq"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

type MarketPlaceStat struct {
	LastIndex uint64 `json:"lastIndex"`
	ID        uint64 `json:"id" gorm:"primaryKey"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime:milli"`  // Set to current unix seconds on updaing or if it is zero on creating
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime:milli"` // Use unix seconds as creating time
}

type MarketPlaceNFT struct {
	Identifier           string         `json:"identifier"`
	Collection           string         `json:"collection"`
	Nonce                uint64         `json:"nonce"`
	NFTType              string         `json:"type"`
	Creator              string         `json:"creator"`
	Royalties            uint64         `json:"royalties"`
	URIs                 pq.StringArray `json:"uris"`
	URL                  string         `json:"url"`
	IsWhitelistedStorage bool           `json:"isWhitelistedStorage"`
	Metadata             JSONB          `json:"metadata"`
	Ticker               string         `json:"ticker"`
}
