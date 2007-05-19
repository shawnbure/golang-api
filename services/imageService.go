package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	images2 "github.com/erdsea/erdsea-api/data/entities/images"

	"github.com/erdsea/erdsea-api/storage"
)

var (
	maxProfileImageSize = 512 * 1024
	maxCoverImageSize   = 1024 * 1024
	errorProfileTooBig  = errors.New(fmt.Sprintf("profile image exceeded max size of %d", maxProfileImageSize))
	errorCoverTooBig    = errors.New(fmt.Sprintf("profile image exceeded max size of %d", maxCoverImageSize))
)

func SetAccountProfileImage(accountId uint64, image *string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxProfileImageSize {
		return errorProfileTooBig
	}

	profileImage := images2.AccountProfileImage{
		ImageBase64: *image,
		AccountID:   accountId,
	}
	return storage.AddOrUpdateAccountProfileImage(&profileImage)
}

func SetAccountCoverImage(accountId uint64, image *string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxCoverImageSize {
		return errorCoverTooBig
	}

	coverImage := images2.AccountCoverImage{
		ImageBase64: *image,
		AccountID:   accountId,
	}

	return storage.AddOrUpdateAccountCoverImage(&coverImage)
}

func SetCollectionCoverImage(collectionId uint64, image *string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxCoverImageSize {
		return errorCoverTooBig
	}

	coverImage := images2.CollectionCoverImage{
		ImageBase64:  *image,
		CollectionID: collectionId,
	}
	return storage.AddOrUpdateCollectionCoverImage(&coverImage)
}

func SetCollectionProfileImage(collectionId uint64, image *string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxProfileImageSize {
		return errorCoverTooBig
	}

	profileImage := images2.CollectionProfileImage{
		ImageBase64:  *image,
		CollectionID: collectionId,
	}
	return storage.AddOrUpdateCollectionProfileImage(&profileImage)
}

func getByteArrayLenOfBase64EncodedImage(image *string) int {
	return base64.RawStdEncoding.DecodedLen(len(*image))
}
