package services

import (
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
)

type SetWhitelistRequest struct {
	CollectionID uint64 `json:"collectionId"`
	Address      string `json:"address"`
	Amount       uint64 `json:"amount"`
	Type         uint64 `json:"type"`
}

func GetWhitelist(address string) (*entities.Whitelist, error) {
	whitelist, err := storage.GetWhitelistByAddress(address)

	if err != nil {
		return nil, err
	}

	return whitelist, nil
}

func CreateWhitelist(address string, request *SetWhitelistRequest) (*entities.Whitelist, error) {

	whitelist := entities.Whitelist{
		Address:      address,
		CollectionID: request.CollectionID,
		Amount:       request.Amount,
		Type:         request.Type,
		CreatedAt:    uint64(time.Now().Unix()),
		ModifiedAt:   uint64(time.Now().Unix()),
	}

	err := storage.AddWhitelist(&whitelist)
	if err != nil {
		return nil, err
	}

	return &whitelist, err
}

func UpdateWhitelist(whitelist *entities.Whitelist, request *SetWhitelistRequest) error {
	whitelist.CollectionID = request.CollectionID
	whitelist.Address = request.Address
	whitelist.Amount = request.Amount
	whitelist.Type = request.Type
	whitelist.ModifiedAt = uint64(time.Now().Unix())

	return storage.UpdateWhitelist(whitelist)
}
