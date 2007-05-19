package services

import (
	"encoding/hex"
	"fmt"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"math/big"
	"strconv"
	"time"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/interaction"
)

var (
	GetCreatorRoyaltiesView                      = "getCreatorRoyalties"
	GetCreatorLastWithdrawalEpochView            = "getCreatorLastWithdrawalEpoch"
	GetRemainingEpochsUntilClaimView             = "getRemainingEpochsUntilClaim"
	RoyaltiesLocalCacheKeyFormat                 = "Royalties:%s"
	CreatorLastWithdrawalEpochKeyFormat          = "CLWE:%s"
	CreatorRemainingEpochsUntilWithdrawKeyFormat = "CWRE:%s"
	RoyaltiesExpirePeriod                        = 15 * time.Minute
	CreatorLastWithdrawalExpirePeriod            = 15 * time.Minute
	RemainingEpochsUntilClaimExpirePeriod        = 15 * time.Minute
)

func GetCreatorRoyalties(marketplaceAddress string, address string) (float64, error) {
	localCacher := cache.GetLocalCacher()
	key := fmt.Sprintf(RoyaltiesLocalCacheKeyFormat, address)

	priceVal, errRead := localCacher.Get(key)
	if errRead == nil {
		return priceVal.(float64), nil
	}

	amountMaybeEmpty, err := DoGetCreatorRoyaltiesVmQuery(marketplaceAddress, address)
	if err != nil {
		return 0, err
	}

	amount := "00"
	if len(amountMaybeEmpty) != 0 {
		amount = amountMaybeEmpty
	}

	amountNominal, err := GetPriceNominal(amount)
	if err != nil {
		log.Debug("could not get amount nominal")
		return 0, err
	}

	errSet := localCacher.SetWithTTL(key, amountNominal, RoyaltiesExpirePeriod)
	if errSet != nil {
		log.Debug("could not cache result", errSet)
	}

	return amountNominal, nil
}

func GetCreatorLastWithdrawalEpoch(marketplaceAddress string, address string) (int64, error) {
	localCacher := cache.GetLocalCacher()
	key := fmt.Sprintf(CreatorLastWithdrawalEpochKeyFormat, address)

	priceVal, errRead := localCacher.Get(key)
	if errRead == nil {
		return priceVal.(int64), nil
	}

	epochsHexMaybeEmpty, err := DoGetCreatorLastWithdrawalEpochVmQuery(marketplaceAddress, address)
	if err != nil {
		return 0, err
	}

	epochsHex := "00"
	if len(epochsHexMaybeEmpty) != 0 {
		epochsHex = epochsHexMaybeEmpty
	}

	epochs, err := strconv.ParseInt(epochsHex, 16, 0)
	if err != nil {
		log.Debug("could not decode epochs hex")
		return 0, err
	}

	errSet := localCacher.SetWithTTL(key, epochs, CreatorLastWithdrawalExpirePeriod)
	if errSet != nil {
		log.Debug("could not cache result", errSet)
	}

	return epochs, nil
}

func GetCreatorRemainingEpochsUntilWithdraw(marketplaceAddress string, address string) (int64, error) {
	localCacher := cache.GetLocalCacher()
	key := fmt.Sprintf(CreatorRemainingEpochsUntilWithdrawKeyFormat, address)

	priceVal, errRead := localCacher.Get(key)
	if errRead == nil {
		return priceVal.(int64), nil
	}

	epochsHexMaybeEmpty, err := DoGetRemainingEpochsUntilClaim(marketplaceAddress, address)
	if err != nil {
		return 0, err
	}

	epochsHex := "00"
	if len(epochsHexMaybeEmpty) != 0 {
		epochsHex = epochsHexMaybeEmpty
	}

	epochs, err := strconv.ParseInt(epochsHex, 16, 0)
	if err != nil {
		log.Debug("could not decode epochs hex")
		return 0, err
	}

	errSet := localCacher.SetWithTTL(key, epochs, RemainingEpochsUntilClaimExpirePeriod)
	if errSet != nil {
		log.Debug("could not cache result", errSet)
	}

	return epochs, nil
}

func DoGetCreatorRoyaltiesVmQuery(marketplaceAddress string, address string) (string, error) {
	bi := interaction.GetBlockchainInteractor()

	addressDecoded, err := data.NewAddressFromBech32String(address)
	if err != nil {
		return "", err
	}

	addressHex := hex.EncodeToString(addressDecoded.AddressBytes())
	result, err := bi.DoVmQuery(marketplaceAddress, GetCreatorRoyaltiesView, []string{addressHex})
	if err != nil || len(result) == 0 {
		return "", nil
	}

	deposit := big.NewInt(0).SetBytes(result[0])
	depositBytes := deposit.Bytes()
	return hex.EncodeToString(depositBytes), nil
}

func DoGetCreatorLastWithdrawalEpochVmQuery(marketplaceAddress string, address string) (string, error) {
	bi := interaction.GetBlockchainInteractor()

	addressDecoded, err := data.NewAddressFromBech32String(address)
	if err != nil {
		return "", err
	}

	addressHex := hex.EncodeToString(addressDecoded.AddressBytes())
	result, err := bi.DoVmQuery(marketplaceAddress, GetCreatorLastWithdrawalEpochView, []string{addressHex})
	if err != nil || len(result) == 0 {
		return "", nil
	}

	deposit := big.NewInt(0).SetBytes(result[0])
	depositBytes := deposit.Bytes()
	return hex.EncodeToString(depositBytes), nil
}

func DoGetRemainingEpochsUntilClaim(marketplaceAddress string, address string) (string, error) {
	bi := interaction.GetBlockchainInteractor()

	addressDecoded, err := data.NewAddressFromBech32String(address)
	if err != nil {
		return "", err
	}

	addressHex := hex.EncodeToString(addressDecoded.AddressBytes())
	result, err := bi.DoVmQuery(marketplaceAddress, GetRemainingEpochsUntilClaimView, []string{addressHex})
	if err != nil || len(result) == 0 {
		return "", nil
	}

	deposit := big.NewInt(0).SetBytes(result[0])
	depositBytes := deposit.Bytes()
	return hex.EncodeToString(depositBytes), nil
}
