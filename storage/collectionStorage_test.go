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
	error := addNewCollection(&collection)
	require.Nil(t, error)

	var collectionRead data.Collection
	txRead := GetDB().Last(&collectionRead)

	require.Nil(t, txRead.Error)
	assert.Equal(t, collectionRead, collection)
}

func Test_GetCollectionById(t *testing.T) {
	connectToDb(t)

	collection, err := getCollectionById(1)
	require.Nil(t, err)
	assert.Equal(t, collection.ID, uint64(1))
}

func Test_GetCollectionsCreatedById(t *testing.T) {
	connectToDb(t)

	collection := defaultCollection()
	error := addNewCollection(&collection)
	require.Nil(t, error)

	otherCollection := defaultCollection()
	error = addNewCollection(&otherCollection)
	require.Nil(t, error)

	collections, error := getCollectionsCreatedBy(0)
	require.Nil(t, error)
	require.GreaterOrEqual(t, len(collections), 2)
}

func Test_GetCollectionByName(t *testing.T) {
	connectToDb(t)

	collectionName := "insane_unique_name"
	collection := defaultCollection()
	collection.Name = collectionName
	error := addNewCollection(&collection)
	require.Nil(t, error)

	retrievedCollection, error := getCollectionByName(collectionName)
	require.Nil(t, error)
	require.GreaterOrEqual(t, retrievedCollection.Name, collectionName)
}
