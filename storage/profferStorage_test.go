package storage

import (
	"testing"

	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/stretchr/testify/require"
)

func Test_DeleteProffer(t *testing.T) {
	connectToTestDb()

	proffer := entities.Proffer{
		Type:          "Offer",
		AmountNominal: 1,
		TokenID:       1,
		OfferorID:     1,
	}
	err := AddProffer(&proffer)
	require.Nil(t, err)

	err = DeleteOfferByTokenIdAndAccountId(1, 1)
	require.Nil(t, err)
}
