package process

import (
	"encoding/json"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/services"
)

var log = logger.GetOrCreate("eventProcessor")

const (
	putNFTForSaleEventName = "put_nft_for_sale"
	buyNFTEventName        = "buy_nft"
	withdrawNFTEventName   = "withdraw_nft"
	makeOfferEventName     = "make_offer"
	acceptOfferEventName   = "accept_offer"
	startAuctionEventName  = "start_auction"
	placeBidEventName      = "place_bid"
	endAuctionEventName    = "end_bid"
	updateDepositEventName = "update_deposit"
	cancelOfferEventName   = "cancel_offer"

	saveEventsTTL = 5 * time.Minute
)

type EventProcessor struct {
	addressSet     map[string]bool
	identifiersSet map[string]bool

	eventsPool chan []entities.Event

	localCacher *cache.LocalCacher
	monitor     *observerMonitor
}

func NewEventProcessor(
	addresses []string,
	identifiers []string,
	monitor *observerMonitor,
) *EventProcessor {
	addrSet := map[string]bool{}
	idSet := map[string]bool{}

	for _, addr := range addresses {
		addrSet[addr] = true
	}

	for _, id := range identifiers {
		idSet[id] = true
	}

	processor := &EventProcessor{
		addressSet:     addrSet,
		identifiersSet: idSet,
		eventsPool:     make(chan []entities.Event),
		localCacher:    cache.GetLocalCacher(),
		monitor:        monitor,
	}

	go processor.PoolWorker()

	return processor
}

func (e *EventProcessor) PoolWorker() {
	for eventArray := range e.eventsPool {
		for _, event := range eventArray {
			if len(event.Topics) == 0 {
				continue
			}

			switch getEventName(&event) {
			case putNFTForSaleEventName:
				e.onEventPutNftForSale(event)
			case buyNFTEventName:
				e.onEventBuyNft(event)
			case withdrawNFTEventName:
				e.onEventWithdrawNft(event)
			case makeOfferEventName:
				e.onEventMakeOffer(event)
			case acceptOfferEventName:
				e.onEventAcceptOffer(event)
			case startAuctionEventName:
				e.onEventStartAuction(event)
			case placeBidEventName:
				e.onEventPlaceBid(event)
			case endAuctionEventName:
				e.onEventEndAuction(event)
			case updateDepositEventName:
				e.onEventUpdateDeposit(event)
			case cancelOfferEventName:
				e.onEventCancelOffer(event)
			}
		}
	}
}

func (e *EventProcessor) OnEvents(blockEvents entities.BlockEvents) {
	var filterableEvents []entities.Event

	for _, event := range blockEvents.Events {
		if e.isEventAccepted(event) {
			filterableEvents = append(filterableEvents, event)
		}
	}

	if len(filterableEvents) > 0 {
		err := e.localCacher.SetWithTTLSync(blockEvents.Hash, filterableEvents, saveEventsTTL)
		if err != nil {
			log.Error(
				"could not store events at block",
				"headerHash", blockEvents.Hash,
				"err", err.Error(),
			)
		}

		log.Info("pushed events to cache for block", "headerHash", blockEvents.Hash)
	}
}

func (e *EventProcessor) OnFinalizedEvent(fb entities.FinalizedBlock) {
	if e.monitor.IsEnabled() {
		e.monitor.LivenessChan() <- fb.Hash
	}

	cachedValue, err := e.localCacher.Get(fb.Hash)
	if err != nil {
		log.Warn("could not load events from cache for block",
			"headerHash", fb.Hash,
			"err", err.Error(),
		)
		return
	}

	cachedEvents, ok := cachedValue.([]entities.Event)
	if !ok {
		log.Error(
			"could not cast cached value to []entities.Event",
			"reason", "corrupted cached data",
		)
		return
	}

	if len(cachedEvents) == 0 {
		log.Warn("loaded empty []entities.Event from cache at block", "headerHash", fb.Hash)
		return
	}

	e.eventsPool <- cachedEvents

	err = e.localCacher.Del(fb.Hash)
	if err != nil {
		log.Error("could not delete block events", "err", err.Error())
	}
}

func (e *EventProcessor) isEventAccepted(ev entities.Event) bool {
	return e.addressSet[ev.Address] && e.identifiersSet[ev.Identifier]
}

func (e *EventProcessor) onEventPutNftForSale(event entities.Event) {
	if len(event.Topics) != 13 {
		log.Error("received corrupted putNFTForSale event", "err", "incorrect topics length")
		return
	}

	args := services.ListTokenArgs{
		OwnerAddress:     decodeAddressFromTopic(event.Topics[1]),
		TokenId:          decodeStringFromTopic(event.Topics[2]),
		Nonce:            decodeU64FromTopic(event.Topics[3]),
		TokenName:        decodeStringFromTopic(event.Topics[4]),
		FirstLink:        decodeStringFromTopic(event.Topics[5]),
		SecondLink:       decodeStringFromTopic(event.Topics[6]),
		Hash:             decodeHexStringOrEmptyWhenZeroFromTopic(event.Topics[7]),
		Attributes:       decodeStringFromTopic(event.Topics[8]),
		Price:            decodeBigUintFromTopic(event.Topics[9]),
		RoyaltiesPercent: decodeU64FromTopic(event.Topics[10]),
		Timestamp:        decodeU64FromTopic(event.Topics[11]),
		TxHash:           decodeTxHashFromTopic(event.Topics[12]),
	}

	eventJson, err := json.Marshal(args)
	if err == nil {
		log.Debug("onEventPutNftForSale", string(eventJson))
	}

	services.ListToken(args)
}

func (e *EventProcessor) onEventBuyNft(event entities.Event) {
	if len(event.Topics) != 8 {
		log.Error("received corrupted buyNFT event", "err", "incorrect topics length")
		return
	}

	args := services.BuyTokenArgs{
		OwnerAddress: decodeAddressFromTopic(event.Topics[1]),
		BuyerAddress: decodeAddressFromTopic(event.Topics[2]),
		TokenId:      decodeStringFromTopic(event.Topics[3]),
		Nonce:        decodeU64FromTopic(event.Topics[4]),
		Price:        decodeBigUintFromTopic(event.Topics[5]),
		Timestamp:    decodeU64FromTopic(event.Topics[6]),
		TxHash:       decodeTxHashFromTopic(event.Topics[7]),
	}

	eventJson, err := json.Marshal(args)
	if err == nil {
		log.Debug("onEventBuyNft", string(eventJson))
	}

	services.BuyToken(args)
}

func (e *EventProcessor) onEventWithdrawNft(event entities.Event) {
	if len(event.Topics) != 7 {
		log.Error("received corrupted withdrawNFT event", "err", "incorrect topics length")
		return
	}

	args := services.WithdrawTokenArgs{
		OwnerAddress: decodeAddressFromTopic(event.Topics[1]),
		TokenId:      decodeStringFromTopic(event.Topics[2]),
		Nonce:        decodeU64FromTopic(event.Topics[3]),
		Price:        decodeBigUintFromTopic(event.Topics[4]),
		Timestamp:    decodeU64FromTopic(event.Topics[5]),
		TxHash:       decodeTxHashFromTopic(event.Topics[6]),
	}

	eventJson, err := json.Marshal(args)
	if err == nil {
		log.Debug("onEventWithdrawNft", string(eventJson))
	}

	services.WithdrawToken(args)
}

func (e *EventProcessor) onEventMakeOffer(event entities.Event) {
	if len(event.Topics) != 8 {
		log.Error("received corrupted makeOffer event", "err", "incorrect topics length")
		return
	}

	args := services.MakeOfferArgs{
		OfferorAddress: decodeAddressFromTopic(event.Topics[1]),
		TokenId:        decodeStringFromTopic(event.Topics[2]),
		Nonce:          decodeU64FromTopic(event.Topics[3]),
		Amount:         decodeBigUintFromTopic(event.Topics[4]),
		Expire:         decodeU64FromTopic(event.Topics[5]),
		Timestamp:      decodeU64FromTopic(event.Topics[6]),
		TxHash:         decodeTxHashFromTopic(event.Topics[7]),
	}

	eventJson, err := json.Marshal(args)
	if err == nil {
		log.Debug("onEventMakeOffer", string(eventJson))
	}

	_, err = services.MakeOffer(args)
	if err != nil {
		log.Error("could not make offer", err)
	}
}

func (e *EventProcessor) onEventCancelOffer(event entities.Event) {
	if len(event.Topics) != 7 {
		log.Error("received corrupted cancelOffer event", "err", "incorrect topics length")
		return
	}

	args := services.CancelOfferArgs{
		OfferorAddress: decodeAddressFromTopic(event.Topics[1]),
		TokenId:        decodeStringFromTopic(event.Topics[2]),
		Nonce:          decodeU64FromTopic(event.Topics[3]),
		Amount:         decodeBigUintFromTopic(event.Topics[4]),
		Timestamp:      decodeU64FromTopic(event.Topics[5]),
		TxHash:         decodeTxHashFromTopic(event.Topics[6]),
	}

	eventJson, err := json.Marshal(args)
	if err == nil {
		log.Debug("onEventCancelOffer", string(eventJson))
	}

	services.CancelOffer(args)
}

func (e *EventProcessor) onEventAcceptOffer(event entities.Event) {
	if len(event.Topics) != 8 {
		log.Error("received corrupted acceptOffer event", "err", "incorrect topics length")
		return
	}

	args := services.AcceptOfferArgs{
		OwnerAddress:   decodeAddressFromTopic(event.Topics[1]),
		TokenId:        decodeStringFromTopic(event.Topics[2]),
		Nonce:          decodeU64FromTopic(event.Topics[3]),
		OfferorAddress: decodeAddressFromTopic(event.Topics[4]),
		Amount:         decodeBigUintFromTopic(event.Topics[5]),
		Timestamp:      decodeU64FromTopic(event.Topics[6]),
		TxHash:         decodeTxHashFromTopic(event.Topics[7]),
	}

	eventJson, err := json.Marshal(args)
	if err == nil {
		log.Debug("onEventAcceptOffer", string(eventJson))
	}

	services.AcceptOffer(args)
}

func (e *EventProcessor) onEventStartAuction(event entities.Event) {
	if len(event.Topics) != 15 {
		log.Error("received corrupted startAuction event", "err", "incorrect topics length")
		return
	}

	args := services.StartAuctionArgs{
		OwnerAddress:     decodeAddressFromTopic(event.Topics[1]),
		TokenId:          decodeStringFromTopic(event.Topics[2]),
		Nonce:            decodeU64FromTopic(event.Topics[3]),
		TokenName:        decodeStringFromTopic(event.Topics[4]),
		FirstLink:        decodeStringFromTopic(event.Topics[5]),
		SecondLink:       decodeStringFromTopic(event.Topics[6]),
		Hash:             decodeHexStringOrEmptyWhenZeroFromTopic(event.Topics[7]),
		Attributes:       decodeStringFromTopic(event.Topics[8]),
		MinBid:           decodeBigUintFromTopic(event.Topics[9]),
		StartTime:        decodeU64FromTopic(event.Topics[10]),
		Deadline:         decodeU64FromTopic(event.Topics[11]),
		RoyaltiesPercent: decodeU64FromTopic(event.Topics[12]),
		Timestamp:        decodeU64FromTopic(event.Topics[13]),
		TxHash:           decodeTxHashFromTopic(event.Topics[14]),
	}

	eventJson, err := json.Marshal(args)
	if err == nil {
		log.Debug("onEventAcceptOffer", string(eventJson))
	}

	_, err = services.StartAuction(args)
	if err != nil {
		log.Error("could not start auction", err)
	}
}

func (e *EventProcessor) onEventPlaceBid(event entities.Event) {
	if len(event.Topics) != 7 {
		log.Error("received corrupted placeBid event", "err", "incorrect topics length")
		return
	}

	args := services.PlaceBidArgs{
		Offeror:   decodeAddressFromTopic(event.Topics[1]),
		TokenId:   decodeStringFromTopic(event.Topics[2]),
		Nonce:     decodeU64FromTopic(event.Topics[3]),
		Amount:    decodeBigUintFromTopic(event.Topics[4]),
		Timestamp: decodeU64FromTopic(event.Topics[5]),
		TxHash:    decodeTxHashFromTopic(event.Topics[6]),
	}

	eventJson, err := json.Marshal(args)
	if err == nil {
		log.Debug("onEventPlaceBid", string(eventJson))
	}

	_, err = services.PlaceBid(args)
	if err != nil {
		log.Error("could not place bid", err)
	}
}

func (e *EventProcessor) onEventEndAuction(event entities.Event) {
	if len(event.Topics) != 8 {
		log.Error("received corrupted makeOffer event", "err", "incorrect topics length")
		return
	}

	args := services.EndAuctionArgs{
		TokenId:   decodeStringFromTopic(event.Topics[2]),
		Nonce:     decodeU64FromTopic(event.Topics[3]),
		Winner:    decodeAddressFromTopic(event.Topics[4]),
		Amount:    decodeBigUintFromTopic(event.Topics[5]),
		Timestamp: decodeU64FromTopic(event.Topics[6]),
		TxHash:    decodeTxHashFromTopic(event.Topics[7]),
	}

	eventJson, err := json.Marshal(args)
	if err == nil {
		log.Debug("onEventPlaceBid", string(eventJson))
	}

	services.EndAuction(args)
}

func (e *EventProcessor) onEventUpdateDeposit(event entities.Event) {
	if len(event.Topics) != 3 {
		log.Error("received corrupted makeOffer event", "err", "incorrect topics length")
		return
	}

	args := services.DepositUpdateArgs{
		Owner:  decodeAddressFromTopic(event.Topics[1]),
		Amount: decodeBigUintFromTopic(event.Topics[2]),
	}

	eventJson, err := json.Marshal(args)
	if err == nil {
		log.Debug("onEventUpdateDeposit", string(eventJson))
	}

	err = services.UpdateDeposit(args)
	if err != nil {
		log.Error("could not upgrade deposit", err)
	}
}

func getEventName(event *entities.Event) string {
	return string(event.Topics[0])
}
