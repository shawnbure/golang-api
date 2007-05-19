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
	GetDepositView             = "getEgldDeposit"
	DepositLocalCacheKeyFormat = "Deposit:%s"
	DepositExpirePeriod        = 15 * time.Minute
)

func UpdateDeposit(args DepositUpdateArgs) error {
	if args.Owner == ZeroAddress {
		return nil
	}

	localCacher := cache.GetLocalCacher()
	key := fmt.Sprintf(DepositLocalCacheKeyFormat, args.Owner)

	if len(args.Amount) == 0 {
		args.Amount = "00"
	}

	_, err := GetOrAddAccountCacheInfo(args.Owner)
	if err != nil {
		_, innerErr := GetOrCreateAccount(args.Owner)
		if innerErr != nil {
			log.Debug("cannot create account", innerErr)
		}
	}

	depositNominal, err := GetPriceNominal(args.Amount)
	if err != nil {
		log.Debug("could not get price nominal")
		return nil
	}

	errSet := localCacher.SetWithTTLSync(key, depositNominal, DepositExpirePeriod)
	if errSet != nil {
		log.Debug("could not set cache", errSet)
		return nil
	}

	return nil
}

func GetDeposit(marketplaceAddress string, address string) (float64, error) {
	localCacher := cache.GetLocalCacher()
	key := fmt.Sprintf(DepositLocalCacheKeyFormat, address)

	priceVal, errRead := localCacher.Get(key)
	if errRead == nil {
		return priceVal.(float64), nil
	}

	deposit, err := DoGetDepositVmQuery(marketplaceAddress, address)
	if err != nil {
		return 0, err
	}

	depositNominal, err := GetPriceNominal(deposit)
	if err != nil {
		log.Debug("could not get price nominal")
		return 0, err
	}

	errSet := localCacher.SetWithTTLSync(key, depositNominal, DepositExpirePeriod)
	if errSet != nil {
		log.Debug("could not cache result", errSet)
	}

	return depositNominal, nil
}

func DoGetDepositVmQuery(marketplaceAddress string, address string) (string, error) {
	bi := interaction.GetBlockchainInteractor()

	result, err := bi.DoVmQuery(marketplaceAddress, GetDepositView, []string{address})
	if err != nil || len(result) == 0{
		return "", nil
	}

	deposit := big.NewInt(0).SetBytes(result[0])
	depositBytes := deposit.Bytes()
	return hex.EncodeToString(depositBytes), nil
}
