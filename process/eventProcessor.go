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
	putNFTForSaleIdentifier = "putNftForSale"
	buyNFTIdentifier        = "buyNft"
	withdrawNFTIdentifier   = "withdrawNft"

	saveEventsTTL        = 5 * time.Minute
)

type EventProcessor struct {
	addressSet     map[string]bool
	identifiersSet map[string]bool

	eventsPool chan []entities.Event

	localCacher *cache.LocalCacher
}

func NewEventProcessor(
	addresses []string,
	identifiers []string,
	localCacher *cache.LocalCacher,
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
		localCacher:    localCacher,
	}

	go processor.PoolWorker()

	return processor
}

func (e *EventProcessor) PoolWorker() {
	for eventArray := range e.eventsPool {
		for _, event := range eventArray {
			switch event.Identifier {
			case putNFTForSaleIdentifier:
				e.onEventPutNftForSale(event)
			case buyNFTIdentifier:
				e.onEventBuyNft(event)
			case withdrawNFTIdentifier:
				e.onEventWithdrawNft(event)
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

func (e *EventProcessor) OnFinalizedEvent(fb entities.FinalizedBlock) {
	cachedValue, err := e.localCacher.Get(fb.Hash)
	if err != nil {
		log.Error("could not load events from cache for block",
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
		LastLink:         decodeStringFromTopic(event.Topics[6]),
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
