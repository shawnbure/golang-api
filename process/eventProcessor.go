package process

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/services"
)

type EventProcessor struct {
	addressSet     map[string]bool
	identifiersSet map[string]bool
	eventsPool     chan []data.Event
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
		eventsPool:     make(chan []data.Event),
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

func (e *EventProcessor) OnEvents(events []data.Event) {
	var filterableEvents []data.Event

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

func (e *EventProcessor) isEventAccepted(ev data.Event) bool {
	return e.addressSet[ev.Address] && e.identifiersSet[ev.Identifier]
}

func (e *EventProcessor) onEventPutNftForSale(event data.Event) {
	args := services.ListAssetArgs{
		OwnerAddress:   decodeAddressFromTopic(event.Topics[0]),
		TokenId:        decodeStringFromTopic(event.Topics[1]),
		Nonce:          decodeU64FromTopic(event.Topics[2]),
		Uri:            decodeStringFromTopic(event.Topics[3]),
		CollectionName: decodeStringFromTopic(event.Topics[4]),
		Price:          decodeBigUintFromTopic(event.Topics[5]),
		TxHash:         decodeTxHashFromTopic(event.Topics[7]),
	}

	log.Debug("onEventPutNftForSale", args.ToString())
	services.ListAsset(args)
}

func (e *EventProcessor) onEventBuyNft(event data.Event) {
	args := services.BuyAssetArgs{
		OwnerAddress:   decodeAddressFromTopic(event.Topics[0]),
		BuyerAddress:   decodeAddressFromTopic(event.Topics[1]),
		TokenId:        decodeStringFromTopic(event.Topics[2]),
		Nonce:          decodeU64FromTopic(event.Topics[3]),
		Uri:            decodeStringFromTopic(event.Topics[4]),
		CollectionName: decodeStringFromTopic(event.Topics[5]),
		Price:          decodeBigUintFromTopic(event.Topics[6]),
		TxHash:         decodeTxHashFromTopic(event.Topics[8]),
	}

	log.Debug("onEventBuyNft", args.ToString())
	services.BuyAsset(args)
}

func (e *EventProcessor) onEventWithdrawNft(event data.Event) {
	args := services.WithdrawAssetArgs{
		OwnerAddress:   decodeAddressFromTopic(event.Topics[0]),
		TokenId:        decodeStringFromTopic(event.Topics[1]),
		Nonce:          decodeU64FromTopic(event.Topics[2]),
		Uri:            decodeStringFromTopic(event.Topics[3]),
		CollectionName: decodeStringFromTopic(event.Topics[4]),
		Price:          decodeBigUintFromTopic(event.Topics[5]),
		TxHash:         decodeTxHashFromTopic(event.Topics[7]),
	}

	log.Debug("onEventWithdrawNft", args.ToString())
	services.WithdrawAsset(args)
}
