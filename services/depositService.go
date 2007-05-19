package services

import (
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/storage"
)

func UpdateDeposit(args DepositUpdateArgs) (*entities.Deposit, error){
	amountNominal, err := GetPriceNominal(args.Amount)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return nil, err
	}

	accountID := uint64(0)
	accountCacheInfo, err := GetOrAddAccountCacheInfo(args.Owner)
	if err != nil {
		log.Debug("could not get or add acc cache info", err)

		account, innerErr := GetOrCreateAccount(args.Owner)
		if innerErr != nil {
			log.Debug("could not get or add acc", err)
		} else {
			accountID = account.ID
		}
	} else {
		accountID = accountCacheInfo.AccountId
	}

	deposit := entities.Deposit{
		AmountNominal: amountNominal,
		AmountString:  args.Amount,
		OwnerId:       accountID,
	}

	err = storage.UpdateDeposit(&deposit)
	if err != nil {
		innerErr := storage.AddDeposit(&deposit)
		if innerErr != nil {
			log.Debug("could not update or add deposit", innerErr)
			return nil, err
		}
	}

	return &deposit, nil
}

