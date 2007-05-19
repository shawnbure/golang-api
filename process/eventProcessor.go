package process

import (
	"encoding/json"

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
			case "putNftForSale":
				e.onEventPutNftForSale(event)
			case "buyNft":
				e.onEventBuyNft(event)
			case "withdrawNft":
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
	if err != nil {
		log.Debug("onEventPutNftForSale", string(eventJson))
	}

	services.ListToken(args)
}

func (e *EventProcessor) onEventBuyNft(event entities.Event) {
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
	if err != nil {
		log.Debug("onEventBuyNft", string(eventJson))
	}

	services.BuyToken(args)
}

func (e *EventProcessor) onEventWithdrawNft(event entities.Event) {
	args := services.WithdrawTokenArgs{
		OwnerAddress: decodeAddressFromTopic(event.Topics[1]),
		TokenId:      decodeStringFromTopic(event.Topics[2]),
		Nonce:        decodeU64FromTopic(event.Topics[3]),
		Price:        decodeBigUintFromTopic(event.Topics[4]),
		Timestamp:    decodeU64FromTopic(event.Topics[5]),
		TxHash:       decodeTxHashFromTopic(event.Topics[6]),
	}

	eventJson, err := json.Marshal(args)
	if err != nil {
		log.Debug("onEventWithdrawNft", string(eventJson))
	}

	services.WithdrawToken(args)
}
