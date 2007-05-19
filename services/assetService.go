package services

import (
	"fmt"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
)

type ListAssetArgs struct {
	OwnerAddress     string
	TokenId          string
	Nonce            uint64
	Uri              string
	Price            string
	RoyaltiesPercent uint64
	Timestamp        uint64
	TxHash           string
}

var log = logger.GetOrCreate("services")

func (args *ListAssetArgs) ToString() string {
	return fmt.Sprintf(""+
		"OwnerAddress = %s\n"+
		"TokenId = %s\n"+
		"Nonce = %d\n"+
		"Uri = %s\n"+
		"Price = %s\n"+
		"RoyaltiesPercent = %d\n"+
		"Timestamp = %d\n"+
		"TxHash = %s\n",
		args.OwnerAddress,
		args.TokenId,
		args.Nonce,
		args.Uri,
		args.Price,
		args.RoyaltiesPercent,
		args.Timestamp,
		args.TxHash)
}

type BuyAssetArgs struct {
	OwnerAddress string
	BuyerAddress string
	TokenId      string
	Nonce        uint64
	Uri          string
	Price        string
	Timestamp    uint64
	TxHash       string
}

func (args *BuyAssetArgs) ToString() string {
	return fmt.Sprintf(""+
		"OwnerAddress = %s\n"+
		"BuyerAddress = %s\n"+
		"TokenId = %s\n"+
		"Nonce = %d\n"+
		"Uri = %s\n"+
		"Price = %s\n"+
		"Timestamp = %d\n"+
		"TxHash = %s\n",
		args.OwnerAddress,
		args.BuyerAddress,
		args.TokenId,
		args.Nonce,
		args.Uri,
		args.Price,
		args.Timestamp,
		args.TxHash)
}

type WithdrawAssetArgs struct {
	OwnerAddress string
	TokenId      string
	Nonce        uint64
	Uri          string
	Price        string
	Timestamp    uint64
	TxHash       string
}

func (args *WithdrawAssetArgs) ToString() string {
	return fmt.Sprintf(""+
		"OwnerAddress = %s\n"+
		"TokenId = %s\n"+
		"Nonce = %d\n"+
		"Uri = %s\n"+
		"Price = %s\n"+
		"Timestamp = %d\n"+
		"TxHash = %s\n",
		args.OwnerAddress,
		args.TokenId,
		args.Nonce,
		args.Uri,
		args.Price,
		args.Timestamp,
		args.TxHash)
}

func ListAsset(args ListAssetArgs) {
	ownerAccount, err := GetOrCreateAccount(args.OwnerAddress)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}

	collection, err := storage.GetCollectionByTokenId(args.TokenId)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}

	asset := data.Asset{
		TokenID:          args.TokenId,
		Nonce:            args.Nonce,
		Price:            args.Price,
		RoyaltiesPercent: args.RoyaltiesPercent,
		Link:             args.Uri,
		CreatedAt:        args.Timestamp,
		Listed:           true,
		OwnerId:          ownerAccount.ID,
		CollectionID:     collection.ID,
	}

	existingAsset, err := storage.GetAssetByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err == nil {
		asset.ID = existingAsset.ID
		err = storage.UpdateAsset(&asset)
	} else {
		err = storage.AddNewAsset(&asset)
	}
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}

	transaction := data.Transaction{
		Hash:      args.TxHash,
		Type:      "List",
		Price:     args.Price,
		Timestamp: args.Timestamp,
		SellerID:  ownerAccount.ID,
		BuyerID:   0,
		AssetID:   asset.ID,
	}

	err = storage.AddNewTransaction(&transaction)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}
}

func BuyAsset(args BuyAssetArgs) {
	ownerAccount, err := storage.GetAccountByAddress(args.OwnerAddress)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}

	buyerAccount, err := GetOrCreateAccount(args.BuyerAddress)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}

	asset, err := storage.GetAssetByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}

	asset.Listed = false
	asset.OwnerId = 0
	err = storage.UpdateAsset(asset)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}

	transaction := data.Transaction{
		Hash:      args.TxHash,
		Type:      "Buy",
		Price:     args.Price,
		Timestamp: args.Timestamp,
		SellerID:  ownerAccount.ID,
		BuyerID:   buyerAccount.ID,
		AssetID:   asset.ID,
	}

	err = storage.AddNewTransaction(&transaction)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}
}

func WithdrawAsset(args WithdrawAssetArgs) {
	ownerAccount, err := storage.GetAccountByAddress(args.OwnerAddress)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}

	asset, err := storage.GetAssetByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}

	asset.Listed = false
	asset.OwnerId = 0
	err = storage.UpdateAsset(asset)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}

	transaction := data.Transaction{
		Hash:      args.TxHash,
		Type:      "Withdraw",
		Price:     args.Price,
		Timestamp: args.Timestamp,
		SellerID:  0,
		BuyerID:   ownerAccount.ID,
		AssetID:   asset.ID,
	}

	err = storage.AddNewTransaction(&transaction)
	if err != nil {
		log.Debug("Unexpected error: ", err)
		return
	}
}
