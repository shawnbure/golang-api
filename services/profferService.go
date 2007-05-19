package services

import (
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/storage"
)

func MakeOffer(args MakeOfferArgs) (*entities.Proffer, error) {
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

	accountCacheInfo, err := GetOrAddAccountCacheInfo(args.OfferorAddress)
	if err != nil {
		log.Debug("could not get account cache info", err)
		return nil, err
	}

	offer := entities.Proffer{
		Type:          entities.Offer,
		AmountNominal: amountNominal,
		AmountString:  args.Amount,
		Expire:        args.Expire,
		Timestamp:     args.Timestamp,
		TxHash:        args.TxHash,
		TokenID:       tokenCacheInfo.TokenDbId,
		OfferorID:     accountCacheInfo.AccountId,
	}

	err = storage.AddProffer(&offer)
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
	token.Status = entities.None
	token.LastBuyPriceNominal = amountNominal
	err = storage.UpdateToken(token)
	if err != nil {
		log.Debug("could not update token", "err", err)
		return
	}

	err = storage.DeleteProffersForTokenId(token.ID)
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

func PlaceBid(args PlaceBidArgs) (*entities.Proffer, error) {
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

	accountCacheInfo, err := GetOrAddAccountCacheInfo(args.Offeror)
	if err != nil {
		log.Debug("could not get account cache info", err)
		return nil, err
	}

	offer := entities.Proffer{
		Type:          entities.Bid,
		AmountNominal: amountNominal,
		AmountString:  args.Amount,
		Timestamp:     args.Timestamp,
		TxHash:        args.TxHash,
		TokenID:       tokenCacheInfo.TokenDbId,
		OfferorID:     accountCacheInfo.AccountId,
	}

	err = storage.AddProffer(&offer)
	return &offer, nil
}

func CancelOffer(args CancelOfferArgs) {
	amountNominal, err := GetPriceNominal(args.Amount)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return
	}

	tokenCacheInfo, err := GetOrAddTokenCacheInfo(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("could not get token cache info", err)
		return
	}

	accountCacheInfo, err := GetOrAddAccountCacheInfo(args.OfferorAddress)
	if err != nil {
		log.Debug("could not get account cache info", err)
		return
	}

	err = storage.DeleteOffersByTokenIdAccountIdAndAmount(tokenCacheInfo.TokenDbId, accountCacheInfo.AccountId, amountNominal)
	if err != nil {
		log.Debug("could not delete from db", err)
		return
	}
}