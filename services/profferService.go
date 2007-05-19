package services

import (
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/storage"
)

func MakeOffer(args MakeOfferArgs) (*entities.Proffer, error){
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
		Timestamp:     args.Timestamp,
		TxHash:        args.TxHash,
		TokenID:       tokenCacheInfo.TokenDbId,
		OfferorID:     accountCacheInfo.AccountId,
	}

	err = storage.AddProffer(&offer)
	if err != nil {
		log.Debug("could not add offer to db", err)
		return nil, err
	}

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
	token.Listed = false
	token.LastBuyPriceNominal = amountNominal
	err = storage.UpdateToken(token)
	if err != nil {
		log.Debug("could not update token", "err", err)
		return
	}

	err = storage.DeleteProffersForTokenId(token.ID)
	if err != nil {
		log.Debug("could not delete proffers for token", "err", err)
		return
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

func StartAuction(args StartAuctionArgs) {

}

func PlaceBid(args PlaceBidArgs) {

}

func EndAuction(args EndAuctionArgs) {

}
