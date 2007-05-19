package process

import (
	"github.com/erdsea/erdsea-api/data/entities"
	"testing"
)

func TestEventProcessor_OnEvents(t *testing.T) {
	t.Parallel()

	addresses := []string{"erd1", "erd2", "erd3"}
	identifiers := []string{"func1", "func2", "func3"}

	events := []entities.Event{
		{
			Address:    addresses[0],
			Identifier: identifiers[0],
		},
		{
			Address:    addresses[1],
			Identifier: identifiers[1],
		},
		{
			Address:    addresses[2],
			Identifier: "identifiers[2]",
		},
	}

	proc := NewEventProcessor(addresses, identifiers)
	proc.OnEvents(events)
}
