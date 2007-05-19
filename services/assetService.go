package services

import (
	"errors"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
	"math/big"
)

var log = logger.GetOrCreate("services")

func ListAsset(args ListAssetArgs) {
	priceNominal, err := GetPriceNominal(args.Price)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return
	}

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
		PriceNominal:     priceNominal,
		RoyaltiesPercent: GetRoyaltiesPercentNominal(args.RoyaltiesPercent),
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
		Hash:         args.TxHash,
		Type:         data.ListAsset,
		PriceNominal: priceNominal,
		Timestamp:    args.Timestamp,
		SellerID:     ownerAccount.ID,
		BuyerID:      0,
		AssetID:      asset.ID,
	}

	addNewTransaction(&transaction)
}

func BuyAsset(args BuyAssetArgs) {
	priceNominal, err := GetPriceNominal(args.Price)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return
	}

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
	// This was to be reset since the asset will no longer be on the marketplace.
	// Could have been kept like this, but bugs may appear when trying when querying.
	asset.OwnerId = 0
	err = storage.UpdateAsset(asset)
	if err != nil {
		log.Debug("could not update asset", "err", err)
		return
	}

	transaction := data.Transaction{
		Hash:         args.TxHash,
		Type:         data.BuyAsset,
		PriceNominal: priceNominal,
		Timestamp:    args.Timestamp,
		SellerID:     ownerAccount.ID,
		BuyerID:      buyerAccount.ID,
		AssetID:      asset.ID,
	}

	addNewTransaction(&transaction)
}

func WithdrawAsset(args WithdrawAssetArgs) {
	priceNominal, err := GetPriceNominal(args.Price)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return
	}

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
	// This was to be reset since the asset will no longer be on the marketplace.
	// Could have been kept like this, but bugs may appear when trying when querying.
	asset.OwnerId = 0
	err = storage.UpdateAsset(asset)
	if err != nil {
		log.Debug("could not update asset", "err", err)
		return
	}

	transaction := data.Transaction{
		Hash:         args.TxHash,
		Type:         data.WithdrawAsset,
		PriceNominal: priceNominal,
		Timestamp:    args.Timestamp,
		SellerID:     0,
		BuyerID:      ownerAccount.ID,
		AssetID:      asset.ID,
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

func GetPriceNominal(priceHex string) (float64, error) {
	priceBigUint, success := big.NewInt(0).SetString(priceHex, 16)
	if !success {
		return 0, errors.New("could not parse price")
	}

	denominatorBigUint := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(15), nil)
	priceNominalInt := big.NewInt(0).Div(priceBigUint, denominatorBigUint).Int64()
	priceNominal := float64(priceNominalInt) / 1000
	return priceNominal, nil
}

//TODO: Help english please? Denominated or denominal? Goland likes denominated.
func GetPriceDenominated(price float64) *big.Int {
	priceUint := int64(price * 1000)
	denominatorBigUint := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(15), nil)

	priceBigUint := big.NewInt(0).Mul(big.NewInt(priceUint), denominatorBigUint)
	return priceBigUint
}

func GetRoyaltiesPercentNominal(percent uint64) float64 {
	return float64(percent) / 1000
}
