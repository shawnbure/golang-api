package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/erdsea/erdsea-api/cdn"
	"github.com/erdsea/erdsea-api/storage"
)

var (
	maxProfileImageSize = 512 * 1024
	maxCoverImageSize   = 1024 * 1024
	errorProfileTooBig  = errors.New(fmt.Sprintf("profile image exceeded max size of %d", maxProfileImageSize))
	errorCoverTooBig    = errors.New(fmt.Sprintf("profile image exceeded max size of %d", maxCoverImageSize))
	CoverSuffix         = ".cover"
	ProfileSuffix       = ".profile"
	ctx                 = context.Background()
)

func SetAccountProfileImage(accountAddress string, accountId uint64, image *string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxProfileImageSize {
		return errorProfileTooBig
	}

	imgId := accountAddress + ProfileSuffix
	response, err := cdn.UploadToCloudy(ctx, *image, imgId)
	if err != nil {
		return err
	}

	err = storage.UpdateAccountProfileWhereId(accountId, response.SecureURL)
	if err != nil {
		return err
	}

	return nil
}

func SetAccountCoverImage(accountAddress string, accountId uint64, image *string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxCoverImageSize {
		return errorCoverTooBig
	}

	imgId := accountAddress + CoverSuffix
	response, err := cdn.UploadToCloudy(ctx, *image, imgId)
	if err != nil {
		return err
	}

	err = storage.UpdateAccountCoverWhereId(accountId, response.SecureURL)
	if err != nil {
		return err
	}

	return nil
}

func SetCollectionCoverImage(tokenId string, collectionId uint64, image *string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxCoverImageSize {
		return errorCoverTooBig
	}

	imgId := tokenId + CoverSuffix
	response, err := cdn.UploadToCloudy(ctx, *image, imgId)
	if err != nil {
		return err
	}

	err = storage.UpdateCollectionCoverWhereId(collectionId, response.SecureURL)
	if err != nil {
		return err
	}

	return nil
}

func SetCollectionProfileImage(tokenId string, collectionId uint64, image *string) error {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize > maxProfileImageSize {
		return errorCoverTooBig
	}

	imgId := tokenId + ProfileSuffix
	response, err := cdn.UploadToCloudy(ctx, *image, imgId)
	if err != nil {
		return err
	}

	err = storage.UpdateCollectionProfileWhereId(collectionId, response.SecureURL)
	if err != nil {
		return err
	}

	return nil
}

func getByteArrayLenOfBase64EncodedImage(image *string) int {
	return base64.StdEncoding.DecodedLen(len(*image))
}
