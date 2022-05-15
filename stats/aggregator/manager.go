package aggregator

import (
	"fmt"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"strconv"
	"sync"
	"time"
)

const StartProjectThreshold = "2022-01-01 00:00:00"
const StartProjectThresholdInt = 2022010100
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

var (
	logInstance = logger.GetOrCreate("aggregator-manager")
)

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
	go m.aggregatedVolumePerHourRunner()
}

func (m *manager) Stop() {
	for _, item := range m.controlChannels {
		item <- true
	}
}

func (m *manager) aggregatedVolumePerHourRunner() {
	m.getAggregatedVolumePerHour()

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
		currentTimeStr := fmt.Sprintf("%04d-%02d-%02d %02d:00:00", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour())

		oneHourBeforeTime := currentTime.Add(-1 * time.Hour)
		oneHourBeforeTimeStr := fmt.Sprintf("%04d-%02d-%02d %02d:00:00", oneHourBeforeTime.Year(), oneHourBeforeTime.Month(), oneHourBeforeTime.Day(), oneHourBeforeTime.Hour())

		intHourStr := fmt.Sprintf("%4d%02d%02d%02d", oneHourBeforeTime.Year(), oneHourBeforeTime.Month(), oneHourBeforeTime.Day(), oneHourBeforeTime.Hour())
		intHour, _ := strconv.ParseInt(intHourStr, 10, 64)

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

			newRecord := entities.AggregatedVolumePerHour{
				Hour:           intHour,
				BuyVolume:      buyVolume,
				ListVolume:     listVolume,
				WithdrawVolume: withdrawVolume,
			}
			err2 := storage.AddOrUpdateAggregatedVolumePerHour(&newRecord)
			if err2 != nil {
				logInstance.Debug(fmt.Sprintf("cannot insert row for %d", intHour))
				break
			}

			index += 1
		} else {
			if intHour > StartProjectThresholdInt {
				_, err := storage.GetOneAggregatedVolumePerHour(intHour)
				if err != nil {
					// insert a new one
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

					newRecord := entities.AggregatedVolumePerHour{
						Hour:           intHour,
						BuyVolume:      buyVolume,
						ListVolume:     listVolume,
						WithdrawVolume: withdrawVolume,
					}
					err2 := storage.AddOrUpdateAggregatedVolumePerHour(&newRecord)
					if err2 != nil {
						logInstance.Debug(fmt.Sprintf("cannot insert row for %d", intHour))
					}
				} else {
					return
				}
			} else {
				return
			}
		}
		subtractIndex++
		time.Sleep(50 * time.Millisecond)
	}
}
