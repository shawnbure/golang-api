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

func (e *EventProcessor) onEventCollectionRegister(event data.Event) {
	ownerAddress := decodeAddressFromTopic(event.Topics[0])
	tokenId := decodeStringFromTopic(event.Topics[1])
	collectionName := decodeStringFromTopic(event.Topics[2])
	collectionDescription := decodeStringFromTopic(event.Topics[3])
	timestamp := decodeU64FromTopic(event.Topics[4])

	fmt.Println(ownerAddress, tokenId, collectionName, collectionDescription, timestamp)
}

func (e *EventProcessor) onEventPutNftForSale(event data.Event) {
	ownerAddress := decodeAddressFromTopic(event.Topics[0])
	tokenId := decodeStringFromTopic(event.Topics[1])
	nonce := decodeU64FromTopic(event.Topics[2])
	uri := decodeStringFromTopic(event.Topics[3])
	collectionName := decodeStringFromTopic(event.Topics[4])
	price := decodeBigUintFromTopic(event.Topics[5])
	timestamp := decodeU64FromTopic(event.Topics[6])

	fmt.Println(ownerAddress, tokenId, nonce, uri, collectionName, price, timestamp)
}

func (e *EventProcessor) onEventBuyNft(event data.Event) {
	ownerAddress := decodeAddressFromTopic(event.Topics[0])
	buyerAddress := decodeAddressFromTopic(event.Topics[1])
	tokenId := decodeStringFromTopic(event.Topics[2])
	nonce := decodeU64FromTopic(event.Topics[3])
	uri := decodeStringFromTopic(event.Topics[4])
	collectionName := decodeStringFromTopic(event.Topics[5])
	price := decodeBigUintFromTopic(event.Topics[6])
	timestamp := decodeU64FromTopic(event.Topics[7])

	fmt.Println(ownerAddress, buyerAddress, tokenId, nonce, uri, collectionName, price, timestamp)
}

func (e *EventProcessor) onEventWithdrawNft(event data.Event) {
	ownerAddress := decodeAddressFromTopic(event.Topics[0])
	tokenId := decodeStringFromTopic(event.Topics[1])
	nonce := decodeU64FromTopic(event.Topics[2])
	uri := decodeStringFromTopic(event.Topics[3])
	collectionName := decodeStringFromTopic(event.Topics[4])
	price := decodeBigUintFromTopic(event.Topics[5])
	timestamp := decodeU64FromTopic(event.Topics[6])

	fmt.Println(ownerAddress, tokenId, nonce, uri, collectionName, price, timestamp)
}
