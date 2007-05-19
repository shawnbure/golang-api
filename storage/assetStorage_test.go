package storage

import (
	"github.com/erdsea/erdsea-api/data"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AddNewAsset(t *testing.T) {
	connectToTestDb()

	asset := defaultAsset()
	err := AddNewAsset(&asset)
	require.Nil(t, err)

	var assetRead data.Asset
	txRead := GetDB().Last(&assetRead)

	require.Nil(t, txRead.Error)
	require.Equal(t, assetRead, asset)
}

func Test_UpdateAsset(t *testing.T) {
	connectToTestDb()

	asset := defaultAsset()
	err := AddNewAsset(&asset)
	require.Nil(t, err)

	asset.TokenID = "new_token_id"
	err = UpdateAsset(&asset)

	var assetRead data.Asset
	txRead := GetDB().Last(&assetRead)

	require.Nil(t, txRead.Error)
	require.Equal(t, assetRead, asset)
}

func Test_GetAssetById(t *testing.T) {
	connectToTestDb()

	asset := defaultAsset()
	err := AddNewAsset(&asset)
	require.Nil(t, err)

	assetRead, err := GetAssetById(asset.ID)
	require.Nil(t, err)
	require.Equal(t, assetRead, &asset)
}

func Test_GetAssetByTokenIdAndNonce(t *testing.T) {
	connectToTestDb()

	asset := defaultAsset()
	asset.TokenID = "unique_token_id"
	asset.Nonce = uint64(100)

	err := AddNewAsset(&asset)
	require.Nil(t, err)

	assetRead, err := GetAssetByTokenIdAndNonce(asset.TokenID, asset.Nonce)
	require.Nil(t, err)
	require.Equal(t, assetRead.TokenID, asset.TokenID)
	require.Equal(t, assetRead.Nonce, asset.Nonce)
}

func Test_GetAssetsOwnedBy(t *testing.T) {
	connectToTestDb()
	ownerId := uint64(1)

	asset := defaultAsset()
	err := AddNewAsset(&asset)
	require.Nil(t, err)

	otherAsset := defaultAsset()
	err = AddNewAsset(&otherAsset)
	require.Nil(t, err)

	assetsRead, err := GetAssetsOwnedBy(ownerId)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(assetsRead), 2)

	for _, assetRead := range assetsRead {
		require.Equal(t, assetRead.OwnerId, ownerId)
	}
}

func Test_GetAssetsByCollectionId(t *testing.T) {
	connectToTestDb()
	collectionId := uint64(1)

	asset := defaultAsset()
	err := AddNewAsset(&asset)
	require.Nil(t, err)

	otherAsset := defaultAsset()
	err = AddNewAsset(&otherAsset)
	require.Nil(t, err)

	assetsRead, err := GetAssetsByCollectionId(collectionId)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(assetsRead), 2)

	for _, assetRead := range assetsRead {
		require.Equal(t, assetRead.CollectionID, collectionId)
	}
}

func Test_CountListedAssetsByCollectionId(t *testing.T) {
	connectToTestDb()

	asset := defaultAsset()
	err := AddNewAsset(&asset)
	require.Nil(t, err)

	otherAsset := defaultAsset()
	err = AddNewAsset(&otherAsset)
	require.Nil(t, err)

	count, err := CountListedAssetsByCollectionId(1)
	require.Nil(t, err)
	require.GreaterOrEqual(t, count, uint64(2))
}

func Test_CountUniqueOwnersWithListedAssetsByCollectionId(t *testing.T) {
	connectToTestDb()

	asset := defaultAsset()
	err := AddNewAsset(&asset)
	require.Nil(t, err)

	otherAsset := defaultAsset()
	err = AddNewAsset(&otherAsset)
	require.Nil(t, err)

	count, err := CountUniqueOwnersWithListedAssetsByCollectionId(1)
	require.Nil(t, err)
	require.Equal(t, count, uint64(1))
}

func defaultAsset() data.Asset {
	return data.Asset{
		TokenID:      "my_token",
		Nonce:        10,
		PriceNominal: 1_000_000_000_000_000_000_000,
		Listed:       true,
		Link:         "link.com",
		OwnerId:      1,
		CollectionID: 1,
	}
}
