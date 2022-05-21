package aggregator

import (
	"fmt"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	StartProjectThreshold    = "2022-01-01 00:00:00"
	StartProjectThresholdInt = 2022010100
	MaxOverComputeThreshold  = 12
	MaxRunnerCount           = 2
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

var (
	logInstance = logger.GetOrCreate("aggregator-manager")
)

// Manager Constructor - It initializes the db configuration params
func (m *manager) init() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.controlChannels = make([]chan bool, MaxRunnerCount)
	for i := 0; i < MaxRunnerCount; i++ {
		m.controlChannels[i] = make(chan bool, 1)
	}
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
	//	go m.aggregatedVolumePerHourRunner()
	go m.aggregatedVolumePerCollectionPerHourRunner()
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

func (m *manager) aggregatedVolumePerCollectionPerHourRunner() {
	m.getAggregatedVolumePerCollectionPerHour()

	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-m.controlChannels[1]:
			ticker.Stop()
			return
		case <-ticker.C:
			m.getAggregatedVolumePerCollectionPerHour()
		}
	}
}

func (m *manager) getAggregatedVolumePerCollectionPerHour() {
	index := 0
	subtractIndex := 0

	// Get all collectionIds from db
	collections, err := storage.GetVerifiedCollections()
	if err != nil {
		return
	}

	collectionIds := []uint64{}
	for _, item := range collections {
		collectionIds = append(collectionIds, item.ID)
	}

	for {
		currentTime := time.Now().UTC().Add(-1 * time.Duration(subtractIndex) * time.Hour)
		currentTimeStr := fmt.Sprintf("%04d-%02d-%02d %02d:00:00", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour())

		oneHourBeforeTime := currentTime.Add(-1 * time.Hour)
		oneHourBeforeTimeStr := fmt.Sprintf("%04d-%02d-%02d %02d:00:00", oneHourBeforeTime.Year(), oneHourBeforeTime.Month(), oneHourBeforeTime.Day(), oneHourBeforeTime.Hour())

		intHourStr := fmt.Sprintf("%4d%02d%02d%02d", oneHourBeforeTime.Year(), oneHourBeforeTime.Month(), oneHourBeforeTime.Day(), oneHourBeforeTime.Hour())
		intHour, _ := strconv.ParseInt(intHourStr, 10, 64)

		type tempVolumeStruct struct {
			BuyVolume      float64
			ListVolume     float64
			WithdrawVolume float64
		}

		if index < MaxOverComputeThreshold {
			records, err := storage.GetAggregatedTradedVolumePerCollectionHourly(oneHourBeforeTimeStr, currentTimeStr)
			if err != nil {
				break
			}

			tempList := make(map[uint64]tempVolumeStruct)
			tempIds := []uint64{}

			for _, item := range records {
				p := tempList[item.CollectionId]
				switch item.Type {
				case string(entities.BuyToken):
					p.BuyVolume = item.Total
				case string(entities.ListToken):
					p.ListVolume = item.Total
				case string(entities.WithdrawToken):
					p.WithdrawVolume = item.Total
				}
				tempList[item.CollectionId] = p
				tempIds = append(tempIds, item.CollectionId)
			}

			for key, item := range tempList {
				newRecord := entities.AggregatedVolumePerCollectionPerHour{
					Hour:           intHour,
					BuyVolume:      item.BuyVolume,
					ListVolume:     item.ListVolume,
					WithdrawVolume: item.WithdrawVolume,
					CollectionId:   key,
				}
				err2 := storage.AddOrUpdateAggregatedVolumePerCollectionPerHour(&newRecord)
				if err2 != nil {
					logInstance.Debug(fmt.Sprintf("cannot insert row for %d and collection %d", intHour, key))
				}
			}

			// Get collections that does not exist here
			for _, id := range collectionIds {
				indexT := sort.Search(len(tempIds), func(i int) bool {
					return id == tempIds[i]
				})

				if indexT <= 0 || indexT > len(tempIds) {
					newRecord := entities.AggregatedVolumePerCollectionPerHour{
						Hour:           intHour,
						BuyVolume:      0.0,
						ListVolume:     0.0,
						WithdrawVolume: 0.0,
						CollectionId:   id,
					}
					err2 := storage.AddOrUpdateAggregatedVolumePerCollectionPerHour(&newRecord)
					if err2 != nil {
						logInstance.Debug(fmt.Sprintf("cannot insert row for %d and collection %d", intHour, id))
					}
				}
			}

			index += 1
		} else {
			if intHour > StartProjectThresholdInt {
				fetchedRecords, _ := storage.GetOneAggregatedVolumePerCollectionPerHour(intHour)

				if len(fetchedRecords) == 0 {
					records, err := storage.GetAggregatedTradedVolumePerCollectionHourly(oneHourBeforeTimeStr, currentTimeStr)
					if err != nil {
						break
					}

					tempList := make(map[uint64]tempVolumeStruct)
					tempIds := []uint64{}

					for _, item := range records {
						p := tempList[item.CollectionId]
						switch item.Type {
						case string(entities.BuyToken):
							p.BuyVolume = item.Total
						case string(entities.ListToken):
							p.ListVolume = item.Total
						case string(entities.WithdrawToken):
							p.WithdrawVolume = item.Total
						}
						tempList[item.CollectionId] = p
						tempIds = append(tempIds, item.CollectionId)
					}

					for key, item := range tempList {
						newRecord := entities.AggregatedVolumePerCollectionPerHour{
							Hour:           intHour,
							BuyVolume:      item.BuyVolume,
							ListVolume:     item.ListVolume,
							WithdrawVolume: item.WithdrawVolume,
							CollectionId:   key,
						}
						err2 := storage.AddOrUpdateAggregatedVolumePerCollectionPerHour(&newRecord)
						if err2 != nil {
							logInstance.Debug(fmt.Sprintf("cannot insert row for %d and collection %d", intHour, key))
							break
						}
					}

					// Get collections that does not exist here
					for _, id := range collectionIds {
						indexT := sort.Search(len(tempIds), func(i int) bool {
							return id == tempIds[i]
						})

						if indexT <= 0 || indexT > len(tempIds) {
							newRecord := entities.AggregatedVolumePerCollectionPerHour{
								Hour:           intHour,
								BuyVolume:      0.0,
								ListVolume:     0.0,
								WithdrawVolume: 0.0,
								CollectionId:   id,
							}
							err2 := storage.AddOrUpdateAggregatedVolumePerCollectionPerHour(&newRecord)
							if err2 != nil {
								logInstance.Debug(fmt.Sprintf("cannot insert row for %d and collection %d", intHour, id))
							}
						}
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
