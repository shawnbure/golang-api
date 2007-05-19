package services

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
)

var log = logger.GetOrCreate("services")

func ListAsset(args ListAssetArgs) {
	ownerAccount, err := GetOrCreateAccount(args.OwnerAddress)
	if err != nil {
		log.Debug("could not get or create account", "err", err)
		return
	}

	collection, err := storage.GetCollectionByTokenId(args.TokenId)
	if err != nil {
		log.Debug("could not get collection", "err", err)
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

	var innerErr error
	if err == nil {
		asset.ID = existingAsset.ID
		innerErr = storage.UpdateAsset(&asset)
	} else {
		innerErr = storage.AddNewAsset(&asset)
	}

	if innerErr != nil {
		log.Debug("could not create or update asset", "err", innerErr)
		return
	}

	transaction := data.Transaction{
		Hash:      args.TxHash,
		Type:      data.ListAsset,
		Price:     args.Price,
		Timestamp: args.Timestamp,
		SellerID:  ownerAccount.ID,
		BuyerID:   0,
		AssetID:   asset.ID,
	}

	addNewTransaction(&transaction)
}

func BuyAsset(args BuyAssetArgs) {
	ownerAccount, err := storage.GetAccountByAddress(args.OwnerAddress)
	if err != nil {
		log.Debug("could not get owner account", "err", err)
		return
	}

	buyerAccount, err := GetOrCreateAccount(args.BuyerAddress)
	if err != nil {
		log.Debug("could not get or create account", "err", err)
		return
	}

	asset, err := storage.GetAssetByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("could not get asset", "err", err)
		return
	}

	asset.Listed = false
	// TODO: is this intended ?
	asset.OwnerId = 0
	err = storage.UpdateAsset(asset)
	if err != nil {
		log.Debug("could not update asset", "err", err)
		return
	}

	transaction := data.Transaction{
		Hash:      args.TxHash,
		Type:      data.BuyAsset,
		Price:     args.Price,
		Timestamp: args.Timestamp,
		SellerID:  ownerAccount.ID,
		BuyerID:   buyerAccount.ID,
		AssetID:   asset.ID,
	}

	addNewTransaction(&transaction)
}

func WithdrawAsset(args WithdrawAssetArgs) {
	ownerAccount, err := storage.GetAccountByAddress(args.OwnerAddress)
	if err != nil {
		log.Debug("could not get owner account", err)
		return
	}

	asset, err := storage.GetAssetByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("could not get asset", "err", err)
		return
	}

	asset.Listed = false
	// TODO: is this intended ?
	asset.OwnerId = 0
	err = storage.UpdateAsset(asset)
	if err != nil {
		log.Debug("could not update asset", "err", err)
		return
	}

	transaction := data.Transaction{
		Hash:      args.TxHash,
		Type:      data.WithdrawAsset,
		Price:     args.Price,
		Timestamp: args.Timestamp,
		SellerID:  0,
		BuyerID:   ownerAccount.ID,
		AssetID:   asset.ID,
	}

	addNewTransaction(&transaction)
}

func addNewTransaction(tx *data.Transaction) {
	err := storage.AddNewTransaction(tx)
	if err != nil {
		log.Debug("could not create new transaction", "err", err)
		return
	}
}
