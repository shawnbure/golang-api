package services

import (
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
)

func ListAsset(
	ownerAddress string,
	tokenId string,
	nonce uint64,
	uri string,
	collectionName string,
	price string,
	txHash string) {

	ownerAccount, err := GetOrCreateAccount(ownerAddress)
	if err != nil {
		printError(err)
		return
	}

	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		printError(err)
		return
	}

	asset := data.Asset{
		TokenID:      tokenId,
		Nonce:        nonce,
		Price:        price,
		Link:         uri,
		Listed:       true,
		OwnerId:      ownerAccount.ID,
		CollectionID: collection.ID,
	}

	existingAsset, err := storage.GetAssetByTokenIdAndNonce(tokenId, nonce)
	if err == nil {
		asset.ID = existingAsset.ID
		err = storage.UpdateAsset(&asset)
	} else {
		err = storage.AddNewAsset(&asset)
	}
	if err != nil {
		printError(err)
		return
	}

	transaction := data.Transaction{
		Hash:     txHash,
		Type:     "List",
		Price:    price,
		SellerID: ownerAccount.ID,
		BuyerID:  0,
		AssetID:  asset.ID,
	}

	err = storage.AddNewTransaction(&transaction)
	if err != nil {
		printError(err)
		return
	}
}

func BuyAsset(
	ownerAddress string,
	buyerAddress string,
	tokenId string,
	nonce uint64,
	price string,
	txHash string) {

	ownerAccount, err := storage.GetAccountByAddress(ownerAddress)
	if err != nil {
		printError(err)
		return
	}

	buyerAccount, err := GetOrCreateAccount(buyerAddress)
	if err != nil {
		printError(err)
		return
	}

	asset, err := storage.GetAssetByTokenIdAndNonce(tokenId, nonce)
	if err != nil {
		printError(err)
		return
	}

	asset.Listed = false
	asset.OwnerId = 0
	err = storage.UpdateAsset(asset)
	if err != nil {
		printError(err)
		return
	}

	transaction := data.Transaction{
		Hash:     txHash,
		Type:     "Buy",
		Price:    price,
		SellerID: ownerAccount.ID,
		BuyerID:  buyerAccount.ID,
		AssetID:  asset.ID,
	}

	err = storage.AddNewTransaction(&transaction)
	if err != nil {
		printError(err)
		return
	}
}

func WithdrawAsset(
	ownerAddress string,
	tokenId string,
	nonce uint64,
	price string,
	txHash string) {

	ownerAccount, err := storage.GetAccountByAddress(ownerAddress)
	if err != nil {
		printError(err)
		return
	}

	asset, err := storage.GetAssetByTokenIdAndNonce(tokenId, nonce)
	if err != nil {
		printError(err)
		return
	}

	asset.Listed = false
	asset.OwnerId = 0
	err = storage.UpdateAsset(asset)
	if err != nil {
		printError(err)
		return
	}

	transaction := data.Transaction{
		Hash:     txHash,
		Type:     "Withdraw",
		Price:    price,
		SellerID: 0,
		BuyerID:  ownerAccount.ID,
		AssetID:  asset.ID,
	}

	err = storage.AddNewTransaction(&transaction)
	if err != nil {
		printError(err)
		return
	}
}
