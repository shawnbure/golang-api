package process

import (
	"github.com/erdsea/erdsea-api/data"
)

type EventProcessor struct {
	addressSet     map[string]bool
	identifiersSet map[string]bool
	eventsPool     chan []data.Event

	collected []data.Event
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
		collected:      []data.Event{},
	}

	go processor.PoolWorker()

	return processor
}

func (e *EventProcessor) PoolWorker() {
	events := <-e.eventsPool
	e.collected = append(e.collected, events...)
}

func (e *EventProcessor) OnEvents(events []data.Event) {
	var filterableEvents []data.Event

	for _, event := range events {
		if e.isEventAccepted(event) {
			filterableEvents = append(filterableEvents, event)
		}
	}

	e.eventsPool <- filterableEvents

	return
}

func (e *EventProcessor) isEventAccepted(ev data.Event) bool {
	return e.addressSet[ev.Address] && e.identifiersSet[ev.Identifier]
}
