package process

import (
	"github.com/erdsea/erdsea-api/data"
)

type EventProcessor struct {
	addressSet map[string]bool
	eventsPool chan data.Event

	collected []data.Event
}

func NewEventProcessor(addresses []string) *EventProcessor {
	set := map[string]bool{}

	for _, addr := range addresses {
		set[addr] = true
	}

	processor := &EventProcessor{
		addressSet: set,
		eventsPool: make(chan data.Event, 128),
		collected:  []data.Event{},
	}

	go processor.PoolWorker()

	return processor
}

func (e *EventProcessor) PoolWorker() {
	for event := range e.eventsPool {
		e.collected = append(e.collected, event)
	}
}

func (e *EventProcessor) OnEvents(events []data.Event) {
	for _, event := range events {
		if e.isInSet(event.Address) {
			e.eventsPool <- event
		}
	}

	return
}

func (e *EventProcessor) isInSet(addr string) bool {
	return e.addressSet[addr]
}
