package services

import (
	"encoding/json"
	"errors"
	"gorm.io/datatypes"
	"math/big"
	"strconv"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/storage"
)

type AssetLinkResponse struct {
	Name       string      `json:"name"`
	Image      string      `json:"image"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Value     string `json:"value"`
	TraitType string `json:"trait_type"`
}

var log = logger.GetOrCreate("services")

const (
	minPriceUnit            = 1000
	minPercentUnit          = 1000
	minPercentRoyaltiesUnit = 100
	minPriceDecimals        = 15

	maxAssetLinkResponseSize = 1024
)

var baseExp = big.NewInt(10)

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

	collectionId := uint64(0)
	collection, err := storage.GetCollectionByTokenId(args.TokenId)
	if err == nil {
		collectionId = collection.ID
	}

	asset, err := storage.GetAssetByTokenIdAndNonce(args.TokenId, args.Nonce)

	var innerErr error
	if err != nil {
		newAsset := ConstructNewAssetFromListArgs(args)
		asset = &newAsset
		asset.Listed = true
		asset.PriceNominal = priceNominal
		asset.OwnerId = ownerAccount.ID
		asset.CollectionID = collectionId
		innerErr = storage.AddAsset(asset)
	} else {
		asset.Listed = true
		asset.PriceNominal = priceNominal
		asset.OwnerId = ownerAccount.ID
		asset.CollectionID = collectionId
		innerErr = storage.UpdateAsset(asset)
	}

	if innerErr != nil {
		log.Debug("could not create or update asset", "err", innerErr)
		return
	}

	transaction := entities.Transaction{
		Hash:         args.TxHash,
		Type:         entities.ListAsset,
		PriceNominal: priceNominal,
		Timestamp:    args.Timestamp,
		SellerID:     ownerAccount.ID,
		BuyerID:      0,
		AssetID:      asset.ID,
		CollectionID: collectionId,
	}

	AddTransaction(&transaction)
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

	transaction := entities.Transaction{
		Hash:         args.TxHash,
		Type:         entities.BuyAsset,
		PriceNominal: priceNominal,
		Timestamp:    args.Timestamp,
		SellerID:     ownerAccount.ID,
		BuyerID:      buyerAccount.ID,
		AssetID:      asset.ID,
		CollectionID: asset.CollectionID,
	}

	AddTransaction(&transaction)
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

	transaction := entities.Transaction{
		Hash:         args.TxHash,
		Type:         entities.WithdrawAsset,
		PriceNominal: priceNominal,
		Timestamp:    args.Timestamp,
		SellerID:     0,
		BuyerID:      ownerAccount.ID,
		AssetID:      asset.ID,
		CollectionID: asset.CollectionID,
	}

	AddTransaction(&transaction)
}

func ConstructNewAssetFromListArgs(args ListAssetArgs) entities.Asset {
	asset := entities.Asset{
		TokenID:          args.TokenId,
		Nonce:            args.Nonce,
		RoyaltiesPercent: GetRoyaltiesPercentNominal(args.RoyaltiesPercent),
		MetadataLink:     "",
		CreatedAt:        args.Timestamp,
		Listed:           true,
		Attributes:       datatypes.JSON(""),
		TokenName:        "",
		ImageLink:        "",
		Hash:             args.Hash,
	}

	osResponse, err := GetOSMetadataForAsset(args.LastLink, args.Nonce)
	if err == nil {
		asset.MetadataLink = args.LastLink
		asset.TokenName = osResponse.Name
		asset.ImageLink = osResponse.Image

		attributesJson, innerErr := ConstructAttributesJsonFromResponse(osResponse)
		if innerErr != nil {
			log.Debug("could not parse os response for attributes", "link", args.LastLink)
		} else {
			asset.Attributes = *attributesJson
		}
	} else {
		asset.TokenName = args.TokenName
		asset.ImageLink = args.FirstLink
		asset.Attributes = datatypes.JSON(args.Attributes)
	}

	return asset
}

func GetOSMetadataForAsset(link string, nonce uint64) (*AssetLinkResponse, error) {
	var response AssetLinkResponse

	link = link + "/" + strconv.FormatUint(nonce, 10)
	responseRaw, err := HttpGetRaw(link)
	if err != nil {
		return nil, err
	}
	if len(responseRaw) > maxAssetLinkResponseSize {
		return nil, errors.New("response too long")
	}

	err = json.Unmarshal([]byte(responseRaw), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func ConstructAttributesJsonFromResponse(response *AssetLinkResponse) (*datatypes.JSON, error) {
	attrsMap := make(map[string]string)
	for _, element := range response.Attributes {
		attrsMap[element.TraitType] = element.Value
	}

	attrsBytes, err := json.Marshal(attrsMap)
	if err != nil {
		return nil, err
	}

	attrsJson := datatypes.JSON(attrsBytes)
	return &attrsJson, err
}

func GetPriceNominal(priceHex string) (float64, error) {
	priceBigUint, success := big.NewInt(0).SetString(priceHex, 16)
	if !success {
		return 0, errors.New("could not parse price")
	}

	denominatorBigUint := big.NewInt(0).Exp(baseExp, big.NewInt(minPriceDecimals), nil)
	priceNominalInt := big.NewInt(0).Div(priceBigUint, denominatorBigUint).Int64()
	priceNominal := float64(priceNominalInt) / minPercentUnit
	return priceNominal, nil
}

func GetPriceDenominated(price float64) *big.Int {
	priceInt := int64(price * minPriceUnit)
	if priceInt <= 0 {
		log.Error("price less than min threshold",
			"min_threshold_multiplied", "1",
			"min_threshold_nominal", 1/minPriceUnit,
			"price_int", priceInt,
		)
	}

	denominatorBigUint := big.NewInt(0).Exp(baseExp, big.NewInt(minPriceDecimals), nil)

	priceBigUint := big.NewInt(0).Mul(big.NewInt(priceInt), denominatorBigUint)
	return priceBigUint
}

func GetRoyaltiesPercentNominal(percent uint64) float64 {
	return float64(percent) / minPercentRoyaltiesUnit
}

func AddTransaction(tx *entities.Transaction) {
	err := storage.AddTransaction(tx)
	if err != nil {
		log.Debug("could not create new transaction", "err", err)
		return
	}
}
