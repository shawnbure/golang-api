package process

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/services"
)

type EventProcessor struct {
	addressSet     map[string]bool
	identifiersSet map[string]bool
	eventsPool     chan []entities.Event
}

var log = logger.GetOrCreate("EventProcessor")

func NewEventProcessor(addresses []string, identifiers []string) *EventProcessor {
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
	}

	go processor.PoolWorker()

	return processor
}

func (e *EventProcessor) PoolWorker() {
	for eventArray := range e.eventsPool {
		for _, event := range eventArray {
			switch event.Identifier {
			case "put_nft_for_sale":
				e.onEventPutNftForSale(event)
			case "buy_nft":
				e.onEventBuyNft(event)
			case "withdraw_nft":
				e.onEventWithdrawNft(event)
			}
		}
	}
}

func (e *EventProcessor) OnEvents(events []entities.Event) {
	var filterableEvents []entities.Event

	for _, event := range events {
		if e.isEventAccepted(event) {
			filterableEvents = append(filterableEvents, event)
		}
	}

	if len(filterableEvents) > 0 {
		e.eventsPool <- filterableEvents
	}

	return
}

func (e *EventProcessor) isEventAccepted(ev entities.Event) bool {
	return e.addressSet[ev.Address] && e.identifiersSet[ev.Identifier]
}

func (e *EventProcessor) onEventPutNftForSale(event entities.Event) {
	args := services.ListAssetArgs{
		OwnerAddress:     decodeAddressFromTopic(event.Topics[0]),
		TokenId:          decodeStringFromTopic(event.Topics[1]),
		Nonce:            decodeU64FromTopic(event.Topics[2]),
		Uri:              decodeStringFromTopic(event.Topics[3]),
		Price:            decodeBigUintFromTopic(event.Topics[4]),
		RoyaltiesPercent: decodeU64FromTopic(event.Topics[5]),
		Timestamp:        decodeU64FromTopic(event.Topics[6]),
		TxHash:           decodeTxHashFromTopic(event.Topics[7]),
	}

	log.Debug("onEventPutNftForSale", args.ToString())
	services.ListAsset(args)
}

func (e *EventProcessor) onEventBuyNft(event entities.Event) {
	args := services.BuyAssetArgs{
		OwnerAddress: decodeAddressFromTopic(event.Topics[0]),
		BuyerAddress: decodeAddressFromTopic(event.Topics[1]),
		TokenId:      decodeStringFromTopic(event.Topics[2]),
		Nonce:        decodeU64FromTopic(event.Topics[3]),
		Uri:          decodeStringFromTopic(event.Topics[4]),
		Price:        decodeBigUintFromTopic(event.Topics[5]),
		Timestamp:    decodeU64FromTopic(event.Topics[6]),
		TxHash:       decodeTxHashFromTopic(event.Topics[7]),
	}

	log.Debug("onEventBuyNft", args.ToString())
	services.BuyAsset(args)
}

func (e *EventProcessor) onEventWithdrawNft(event entities.Event) {
	args := services.WithdrawAssetArgs{
		OwnerAddress: decodeAddressFromTopic(event.Topics[0]),
		TokenId:      decodeStringFromTopic(event.Topics[1]),
		Nonce:        decodeU64FromTopic(event.Topics[2]),
		Uri:          decodeStringFromTopic(event.Topics[3]),
		Price:        decodeBigUintFromTopic(event.Topics[4]),
		Timestamp:    decodeU64FromTopic(event.Topics[5]),
		TxHash:       decodeTxHashFromTopic(event.Topics[6]),
	}

	log.Debug("onEventWithdrawNft", args.ToString())
	services.WithdrawAsset(args)
}
