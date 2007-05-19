package process

import (
	"fmt"
	"github.com/erdsea/erdsea-api/data"
)

type EventProcessor struct {
	addressSet     map[string]bool
	identifiersSet map[string]bool
	eventsPool     chan []data.Event
}

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
			case "collection_register":
				e.onEventCollectionRegister(event)
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

func (e* EventProcessor) onEventCollectionRegister(event data.Event) {
	creatorAddress := decodeAddressFromTopic(event.Topics[0])
	tokenId := decodeStringFromTopic(event.Topics[1])
	collectionName := decodeStringFromTopic(event.Topics[2])
	timestamp := decodeU64FromTopic(event.Topics[3])

	fmt.Println(creatorAddress, tokenId, collectionName, timestamp)
}

func (e* EventProcessor) onEventPutNftForSale(event data.Event) {

}

func (e* EventProcessor) onEventBuyNft(event data.Event) {

}

func (e* EventProcessor) onEventWithdrawNft(event data.Event) {

}
