package services

import (
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
)

func MakeOffer(args MakeOfferArgs) (*entities.Offer, error) {
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

	offer := entities.Offer{
		AmountNominal:  amountNominal,
		AmountString:   args.Amount,
		Expire:         args.Expire,
		Timestamp:      args.Timestamp,
		TxHash:         args.TxHash,
		TokenID:        tokenCacheInfo.TokenDbId,
		OfferorAddress: args.OfferorAddress,
	}

	_ = storage.DeleteOfferByOfferorForTokenId(args.OfferorAddress, tokenCacheInfo.TokenDbId)
	err = storage.AddOffer(&offer)
	return &offer, nil
}

func AcceptOffer(args AcceptOfferArgs) {
	amountNominal, err := GetPriceNominal(args.Amount)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return
	}

	buyer, err := GetOrAddAccountCacheInfo(args.OfferorAddress)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return
	}

	token, err := storage.GetTokenByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("could not get token", "err", err)
		return
	}

	sellerId := token.OwnerId
	token.OwnerId = 0
	token.Status = entities.BuyToken
	token.LastBuyPriceNominal = amountNominal
	err = storage.UpdateToken(token)
	if err != nil {
		log.Debug("could not update token", "err", err)
		return
	}

	err = storage.DeleteOffersForTokenId(token.ID)
	if err != nil {
		log.Debug("could not delete proffers for token", "err", err)
	}

	err = storage.DeleteBidsForTokenId(token.ID)
	if err != nil {
		log.Debug("could not delete proffers for token", "err", err)
	}

	transaction := entities.Transaction{
		Hash:         args.TxHash,
		Type:         entities.BuyToken,
		PriceNominal: amountNominal,
		Timestamp:    args.Timestamp,
		SellerID:     sellerId,
		BuyerID:      buyer.AccountId,
		TokenID:      token.ID,
		CollectionID: token.CollectionID,
	}

	AddTransaction(&transaction)
}

func CancelOffer(args CancelOfferArgs) {
	tokenCacheInfo, err := GetOrAddTokenCacheInfo(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("could not get token cache info", err)
		return
	}

	err = storage.DeleteOfferByOfferorForTokenId(args.OfferorAddress, tokenCacheInfo.TokenDbId)
	if err != nil {
		log.Debug("could not delete from db", err)
		return
	}
}

func MakeOfferDtos(offers []entities.Offer) []dtos.OfferDto {
	offerDtos := make([]dtos.OfferDto, len(offers))
	for index := range offers {
		offerorName := ""
		cacheInfo, err := GetOrAddAccountCacheInfo(offers[index].OfferorAddress)
		if err == nil {
			offerorName = cacheInfo.AccountName
		} else {
			log.Debug("cannot get cache info for account", err)
		}

		offerDtos[index] = dtos.OfferDto{
			Offer:       offers[index],
			OfferorName: offerorName,
		}
	}

	return offerDtos
}
