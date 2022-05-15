package storage

import (
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_InsertOrUpdateOneRecord(t *testing.T) {
	connectToTestDb()

	t.Run("Insert a new record", func(t *testing.T) {
		record := defaultAggregatedVolumePerHour()
		err := AddOrUpdateAggregatedVolumePerHour(&record)
		require.Nil(t, err)

		var lastRecordRead entities.AggregatedVolumePerHour
		txRead := GetDB().Last(&lastRecordRead)

		require.Nil(t, txRead.Error)
		require.Equal(t, lastRecordRead.BuyVolume, record.BuyVolume)
	})

	t.Run("Update the existed record", func(t *testing.T) {
		record := defaultAggregatedVolumePerHour()
		record.WithdrawVolume = 2.34
		err := AddOrUpdateAggregatedVolumePerHour(&record)
		require.Nil(t, err)

		getRecord, err := GetOneAggregatedVolumePerHour(record.Hour)
		require.Nil(t, err)

		require.Equal(t, getRecord.WithdrawVolume, record.WithdrawVolume)
	})
}

func defaultAggregatedVolumePerHour() entities.AggregatedVolumePerHour {
	return entities.AggregatedVolumePerHour{
		Hour:           2022051415,
		BuyVolume:      1,
		ListVolume:     2.0,
		WithdrawVolume: 3.2,
	}
}
