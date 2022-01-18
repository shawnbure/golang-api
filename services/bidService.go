package services

import (
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
)

func PlaceBid(args PlaceBidArgs) (*entities.Bid, error) {
	amountNominal, err := GetPriceNominal(args.Amount)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return nil, err
	}

	tokenCacheInfo, err := GetOrAddTokenCacheInfo(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("could not get token cache info", err)
		return nil, err
	}

	bid := entities.Bid{
		BidAmountNominal: amountNominal,
		BidAmountString:  args.Amount,
		Timestamp:        args.Timestamp,
		TxHash:           args.TxHash,
		TokenID:          tokenCacheInfo.TokenDbId,
		BidderAddress:    args.Offeror,
	}

	err = storage.AddBid(&bid)
	return &bid, nil
}

func MakeBidDtos(bids []entities.Bid) []dtos.BidDto {
	bidDtos := make([]dtos.BidDto, len(bids))
	for index := range bids {
		bidderName := ""
		cacheInfo, err := GetOrAddAccountCacheInfo(bids[index].BidderAddress)
		if err == nil {
			bidderName = cacheInfo.AccountName
		} else {
			log.Debug("cannot get cache info for account", err)
		}

		bidDtos[index] = dtos.BidDto{
			Bid:        bids[index],
			BidderName: bidderName,
		}
	}

	return bidDtos
}
