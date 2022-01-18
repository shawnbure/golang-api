package storage

import (
	"testing"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/stretchr/testify/require"
)

func Test_BidOrdering(t *testing.T) {
	connectToTestDb()

	bid1 := entities.Bid{
		TokenID:       10,
		BidderAddress: "erd1",
	}
	err := AddBid(&bid1)
	require.Nil(t, err)

	bid2 := entities.Bid{
		TokenID:       10,
		BidderAddress: "erd2",
	}
	err = AddBid(&bid2)
	require.Nil(t, err)

	bids, err := GetBidsForTokenWithOffsetLimit(10, 0, 1)
	require.Nil(t, err)
	require.Equal(t, 1, len(bids))
	require.Equal(t, "erd2", bids[0].BidderAddress)
}
