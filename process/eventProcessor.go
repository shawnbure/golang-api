package process

import (
	"encoding/json"
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
	args := services.ListTokenArgs{
		OwnerAddress:     decodeAddressFromTopic(event.Topics[0]),
		TokenId:          decodeStringFromTopic(event.Topics[1]),
		Nonce:            decodeU64FromTopic(event.Topics[2]),
		TokenName:        decodeStringFromTopic(event.Topics[3]),
		FirstLink:        decodeStringFromTopic(event.Topics[4]),
		LastLink:         decodeStringFromTopic(event.Topics[5]),
		Hash:             decodeStringFromTopic(event.Topics[6]),
		Attributes:       decodeStringFromTopic(event.Topics[7]),
		Price:            decodeBigUintFromTopic(event.Topics[8]),
		RoyaltiesPercent: decodeU64FromTopic(event.Topics[9]),
		Timestamp:        decodeU64FromTopic(event.Topics[10]),
		TxHash:           decodeTxHashFromTopic(event.Topics[11]),
	}

	eventJson, err := json.Marshal(args)
	if err != nil {
		log.Debug("onEventPutNftForSale", string(eventJson))
	}

	services.ListToken(args)
}

func (e *EventProcessor) onEventBuyNft(event data.Event) {
	args := services.BuyTokenArgs{
		OwnerAddress: decodeAddressFromTopic(event.Topics[0]),
		BuyerAddress: decodeAddressFromTopic(event.Topics[1]),
		TokenId:      decodeStringFromTopic(event.Topics[2]),
		Nonce:        decodeU64FromTopic(event.Topics[3]),
		Price:        decodeBigUintFromTopic(event.Topics[4]),
		Timestamp:    decodeU64FromTopic(event.Topics[5]),
		TxHash:       decodeTxHashFromTopic(event.Topics[6]),
	}

	eventJson, err := json.Marshal(args)
	if err != nil {
		log.Debug("onEventBuyNft", string(eventJson))
	}

	services.BuyToken(args)
}

func (e *EventProcessor) onEventWithdrawNft(event data.Event) {
	args := services.WithdrawTokenArgs{
		OwnerAddress: decodeAddressFromTopic(event.Topics[0]),
		TokenId:      decodeStringFromTopic(event.Topics[1]),
		Nonce:        decodeU64FromTopic(event.Topics[2]),
		Price:        decodeBigUintFromTopic(event.Topics[3]),
		Timestamp:    decodeU64FromTopic(event.Topics[4]),
		TxHash:       decodeTxHashFromTopic(event.Topics[5]),
	}

	eventJson, err := json.Marshal(args)
	if err != nil {
		log.Debug("onEventWithdrawNft", string(eventJson))
	}

	services.WithdrawToken(args)
}
