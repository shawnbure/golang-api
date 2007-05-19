package services

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/erdsea/erdsea-api/data/images"
	"github.com/erdsea/erdsea-api/storage"
)

var (
	maxProfileImageSize = 512 * 1024
	maxCoverImageSize   = 1024 * 1024
	errorProfileTooBig  = errors.New(fmt.Sprintf("profile image exceeded max size of %d", maxProfileImageSize))
	errorCoverTooBig    = errors.New(fmt.Sprintf("profile image exceeded max size of %d", maxCoverImageSize))
	errorNotAuthorized  = errors.New(fmt.Sprintf("user is not authorized"))
)

func GetAccountProfileImage(userAddress string) (*string, error) {
	account, err := storage.GetAccountByAddress(userAddress)
	if err != nil {
		return nil, err
	}

	image, err := storage.GetAccountProfileImageByAccountId(account.ID)
	if err != nil {
		return nil, err
	}

	return &image.ImageBase64, nil
}

func SetAccountProfileImage(userAddress string, image *string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxProfileImageSize {
		return errorProfileTooBig
	}

	account, err := GetOrCreateAccount(userAddress)
	if err != nil {
		return err
	}

	profileImage := images.AccountProfileImage{
		ImageBase64: *image,
		AccountID:   account.ID,
	}
	return storage.AddOrUpdateAccountProfileImage(&profileImage)
}

func GetAccountCoverImage(userAddress string) (*string, error) {
	account, err := storage.GetAccountByAddress(userAddress)
	if err != nil {
		return nil, err
	}

	image, err := storage.GetAccountCoverImageByAccountId(account.ID)
	if err != nil {
		return nil, err
	}

	return &image.ImageBase64, nil
}

func SetAccountCoverImage(userAddress string, image *string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxCoverImageSize {
		return errorCoverTooBig
	}

	account, err := GetOrCreateAccount(userAddress)
	if err != nil {
		return err
	}

	coverImage := images.AccountCoverImage{
		ImageBase64: *image,
		AccountID:   account.ID,
	}

	return storage.AddOrUpdateAccountCoverImage(&coverImage)
}

func GetCollectionCoverImage(collectionName string) (*string, error) {
	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		return nil, err
	}

	image, err := storage.GetCollectionCoverImageByCollectionId(collection.ID)
	if err != nil {
		return nil, err
	}

	return &image.ImageBase64, nil
}

func SetCollectionCoverImage(collectionName string, image *string, authorizedAddress string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxCoverImageSize {
		return errorCoverTooBig
	}

	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		return err
	}

	creatorAccount, err := storage.GetAccountById(collection.CreatorID)
	if err != nil {
		return err
	}
	if creatorAccount.Address != authorizedAddress {
		return errorNotAuthorized
	}

	coverImage := images.CollectionCoverImage{
		ImageBase64:  *image,
		CollectionID: collection.ID,
	}
	return storage.AddOrUpdateCollectionCoverImage(&coverImage)
}

func GetCollectionProfileImage(collectionName string) (*string, error) {
	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		return nil, err
	}

	image, err := storage.GetCollectionProfileImageByCollectionId(collection.ID)
	if err != nil {
		return nil, err
	}

	return &image.ImageBase64, nil
}

func SetCollectionProfileImage(collectionName string, image *string, authorizedAddress string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxProfileImageSize {
		return errorCoverTooBig
	}

	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		return err
	}

	creatorAccount, err := storage.GetAccountById(collection.CreatorID)
	if err != nil {
		return err
	}
	if creatorAccount.Address != authorizedAddress {
		return errorNotAuthorized
	}

	profileImage := images.CollectionProfileImage{
		ImageBase64:  *image,
		CollectionID: collection.ID,
	}
	return storage.AddOrUpdateCollectionProfileImage(&profileImage)
}

func getByteArrayLenOfBase64EncodedImage(image *string) int {
	return base64.RawStdEncoding.DecodedLen(len(*image))
}
