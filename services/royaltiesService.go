package services

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/interaction"
)

var (
	GetCreatorRoyaltiesView             = "getCreatorRoyalties"
	GetCreatorLastWithdrawalEpochView   = "getCreatorLastWithdrawalEpoch"
	RoyaltiesLocalCacheKeyFormat        = "Royalties:%s"
	CreatorLastWithdrawalEpochKeyFormat = "CLWE:%s"
	RoyaltiesExpirePeriod               = 15 * time.Minute
	CreatorLastWithdrawalExpirePeriod   = 15 * time.Minute
)

func GetCreatorRoyalties(marketplaceAddress string, address string) (float64, error) {
	localCacher := cache.GetLocalCacher()
	key := fmt.Sprintf(RoyaltiesLocalCacheKeyFormat, address)

	priceVal, errRead := localCacher.Get(key)
	if errRead == nil {
		return priceVal.(float64), nil
	}

	deposit, err := DoGetCreatorRoyaltiesVmQuery(marketplaceAddress, address)
	if err != nil {
		return 0, err
	}

	depositNominal, err := GetPriceNominal(deposit)
	if err != nil {
		log.Debug("could not get price nominal")
		return 0, err
	}

	errSet := localCacher.SetWithTTL(key, depositNominal, RoyaltiesExpirePeriod)
	if errSet != nil {
		log.Debug("could not cache result", errSet)
	}

	return depositNominal, nil
}

func GetCreatorLastWithdrawalEpoch(marketplaceAddress string, address string) (float64, error) {
	localCacher := cache.GetLocalCacher()
	key := fmt.Sprintf(CreatorLastWithdrawalEpochKeyFormat, address)

	priceVal, errRead := localCacher.Get(key)
	if errRead == nil {
		return priceVal.(float64), nil
	}

	deposit, err := DoGetCreatorLastWithdrawalEpochVmQuery(marketplaceAddress, address)
	if err != nil {
		return 0, err
	}

	depositNominal, err := GetPriceNominal(deposit)
	if err != nil {
		log.Debug("could not get price nominal")
		return 0, err
	}

	errSet := localCacher.SetWithTTL(key, depositNominal, CreatorLastWithdrawalExpirePeriod)
	if errSet != nil {
		log.Debug("could not cache result", errSet)
	}

	return depositNominal, nil
}

func DoGetCreatorRoyaltiesVmQuery(marketplaceAddress string, address string) (string, error) {
	bi := interaction.GetBlockchainInteractor()

	result, err := bi.DoVmQuery(marketplaceAddress, GetCreatorRoyaltiesView, []string{address})
	if err != nil || len(result) == 0 {
		return "", nil
	}

	deposit := big.NewInt(0).SetBytes(result[0])
	depositBytes := deposit.Bytes()
	return hex.EncodeToString(depositBytes), nil
}

func DoGetCreatorLastWithdrawalEpochVmQuery(marketplaceAddress string, address string) (string, error) {
	bi := interaction.GetBlockchainInteractor()

	result, err := bi.DoVmQuery(marketplaceAddress, GetCreatorLastWithdrawalEpochView, []string{address})
	if err != nil || len(result) == 0 {
		return "", nil
	}

	deposit := big.NewInt(0).SetBytes(result[0])
	depositBytes := deposit.Bytes()
	return hex.EncodeToString(depositBytes), nil
}
