package services

import (
	"fmt"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/storage"
)

type CreateNewCollectionArgs struct {
	OwnerAddress          string
	TokenId               string
	CollectionName        string
	CollectionDescription string
}

func (args *CreateNewCollectionArgs) ToString() string {
	return fmt.Sprintf(""+
		"OwnerAddress = %s\n"+
		"TokenId = %s\n"+
		"CollectionName = %s\n"+
		"CollectionDescription = %s\n"+
		args.OwnerAddress,
		args.TokenId,
		args.CollectionName,
		args.CollectionDescription)
}

func CreateNewCollection(args CreateNewCollectionArgs) {
	account, err := GetOrCreateAccount(args.OwnerAddress)
	if err != nil {
		printError(err)
		return
	}

	collection := data.Collection{
		Name:        args.CollectionName,
		TokenID:     args.TokenId,
		Description: args.CollectionDescription,
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
