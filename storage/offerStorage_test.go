package storage

import (
	"testing"

	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/stretchr/testify/require"
)

func Test_DeleteProffer(t *testing.T) {
	connectToTestDb()

	offer := entities.Offer{
		AmountNominal:  1,
		TokenID:        1,
		OfferorAddress: "erd1",
	}
	err := AddOffer(&offer)
	require.Nil(t, err)

	err = DeleteOfferByOfferorForTokenId("erd1", 1)
	require.Nil(t, err)
}
