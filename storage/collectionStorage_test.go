package storage

import (
	"github.com/erdsea/erdsea-api/data"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AddNewCollection(t *testing.T) {
	connectToDb(t)

	collection := defaultCollection()
	err := AddNewCollection(&collection)
	require.Nil(t, err)

	var collectionRead data.Collection
	txRead := GetDB().Last(&collectionRead)

	require.Nil(t, txRead.Error)
	assert.Equal(t, collectionRead, collection)
}

func Test_GetCollectionById(t *testing.T) {
	connectToDb(t)

	collection := defaultCollection()
	err := AddNewCollection(&collection)
	require.Nil(t, err)

	collectionRead, err := GetCollectionById(collection.ID)
	require.Nil(t, err)
	assert.Equal(t, collectionRead.ID, collection.ID)
}

func Test_GetCollectionsCreatedById(t *testing.T) {
	connectToDb(t)

	collection := defaultCollection()
	err := AddNewCollection(&collection)
	require.Nil(t, err)

	otherCollection := defaultCollection()
	err = AddNewCollection(&otherCollection)
	require.Nil(t, err)

	collections, err := GetCollectionsCreatedBy(collection.CreatorID)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(collections), 2)
}

func Test_GetCollectionByName(t *testing.T) {
	connectToDb(t)

	collectionName := "insane_unique_name"
	collection := defaultCollection()
	collection.Name = collectionName
	err := AddNewCollection(&collection)
	require.Nil(t, err)

	retrievedCollection, err := GetCollectionByName(collectionName)
	require.Nil(t, err)
	require.GreaterOrEqual(t, retrievedCollection.Name, collectionName)
}

func defaultCollection() data.Collection {
	return data.Collection{
		Name:      "default",
		TokenID:   "my_token",
		CreatorID: 0,
	}
}
