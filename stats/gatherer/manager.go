package gatherer

import (
	"sync"
)

const (
	MaxRunnerCount = 1
)

// MARK: manager

// Manager object
type manager struct {
	lock            sync.Mutex
	controlChannels []chan bool
}

// MARK: Module variables
var managerInstance *manager = nil
var once sync.Once

// Manager Constructor - It initializes the db configuration params
func (m *manager) init() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.controlChannels = make([]chan bool, MaxRunnerCount)
	for i := 0; i < MaxRunnerCount; i++ {
		m.controlChannels[i] = make(chan bool, 1)
	}
}

// MARK: Public Functions

// GetManager - This function returns singleton instance of Manager
func GetManager() *manager {
	// once used for prevent race condition and manage critical section.
	once.Do(func() {
		managerInstance = &manager{}

		managerInstance.init()
	})
	return managerInstance
}

func (m *manager) Start(blockchainAPI string) {
	// Start hourly aggregator
	go syncRarityRunner(m.controlChannels[0], blockchainAPI)
}

func (m *manager) Stop() {
	for _, item := range m.controlChannels {
		item <- true
	}
}
