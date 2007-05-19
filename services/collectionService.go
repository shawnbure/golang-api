package services

import (
	"fmt"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
)

func CreateNewCollection(ownerAddress string, tokenId string, collectionName string, collectionDescription string) {
	account, err := GetOrCreateAccount(ownerAddress)
	if err != nil {
		printError(err)
		return
	}

	collection := data.Collection{
		Name:        collectionName,
		TokenID:     tokenId,
		Description: collectionDescription,
		CreatorID:   account.ID,
	}

	err = storage.AddNewCollection(&collection)
	if err != nil {
		printError(err)
		return
	}
}

func printError(err error) {
	fmt.Printf("Unexpected error %d\n", err)
}
