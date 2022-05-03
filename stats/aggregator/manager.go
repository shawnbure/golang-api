package aggregator

import (
	"fmt"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	"log"
	"strconv"
	"sync"
	"time"
)

const StartProjectThreshold = "2022-01-01 00:00:00"
const MaxOverComputeThreshold = 12

// MARK: manager

// Manager object
type manager struct {
	lock            sync.Mutex
	controlChannels []chan bool
}

// MARK: Module variables
var managerInstance *manager = nil
var once sync.Once

// Module init function
func init() {
	log.Println("Logger Manager Package Initialized...")
}

// Manager Constructor - It initializes the db configuration params
func (m *manager) init() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.controlChannels = make([]chan bool, 1)
	return
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

func (m *manager) Start() {
	// Start hourly aggregator
	go m.getAggregatedVolumePerHour()
}

func (m *manager) Stop() {
	for _, item := range m.controlChannels {
		item <- true
	}
}

func (m *manager) aggregatedVolumePerHourRunner() {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-m.controlChannels[0]:
			ticker.Stop()
			return
		case <-ticker.C:
			m.getAggregatedVolumePerHour()
		}
	}
}

func (m *manager) getAggregatedVolumePerHour() {
	index := 0
	subtractIndex := 0
	for {
		currentTime := time.Now().UTC().Add(-1 * time.Duration(subtractIndex) * time.Hour)
		currentTimeStr := fmt.Sprintf("%04d-%2d-%2d %2d:00:00", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour())

		oneHourBeforeTime := currentTime.Add(-1 * time.Hour)
		oneHourBeforeTimeStr := fmt.Sprintf("%04d-%2d-%2d %2d:00:00", oneHourBeforeTime.Year(), oneHourBeforeTime.Month(), oneHourBeforeTime.Day(), oneHourBeforeTime.Hour())

		intHourStr := fmt.Sprintf("%4d%2d%2d%2d", oneHourBeforeTime.Year(), oneHourBeforeTime.Month(), oneHourBeforeTime.Day(), oneHourBeforeTime.Hour())
		intHour, _ := strconv.Atoi(intHourStr)

		if index < MaxOverComputeThreshold {
			buyVolume, err := storage.GetAggregatedTradedVolumeHourly(oneHourBeforeTimeStr, currentTimeStr, entities.BuyToken)
			if err != nil {
				break
			}
			listVolume, err := storage.GetAggregatedTradedVolumeHourly(oneHourBeforeTimeStr, currentTimeStr, entities.ListToken)
			if err != nil {
				break
			}
			withdrawVolume, err := storage.GetAggregatedTradedVolumeHourly(oneHourBeforeTimeStr, currentTimeStr, entities.WithdrawToken)
			if err != nil {
				break
			}

			// TODO: insert or update database
			index += 1
		} else {
			// TODO: Check See Whether it exist in database -> If exist return from the function else add it to database (with watching the threshold)
		}
		subtractIndex++
	}
}
