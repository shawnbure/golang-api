package process

import (
	"fmt"
	"github.com/erdsea/erdsea-api/data"
	"testing"
)

func TestEventProcessor_OnEvents(t *testing.T) {
	t.Parallel()

	addresses := []string{"erd1", "erd2", "erd3"}

	events := []data.Event{
		{
			Address:    addresses[0],
			Identifier: "1",
		},
		{
			Address:    addresses[1],
			Identifier: "2",
		},
		{
			Address:    addresses[2],
			Identifier: "3",
		},
	}

	proc := NewEventProcessor(addresses)
	proc.OnEvents(events)

	fmt.Println(proc.collected)
}
