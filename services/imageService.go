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
	errorImageZeroLen   = errors.New(fmt.Sprintf("image length is 0"))
	errorProfileTooBig  = errors.New(fmt.Sprintf("profile image exceeded max size of %d", maxProfileImageSize))
	errorCoverTooBig    = errors.New(fmt.Sprintf("profile image exceeded max size of %d", maxCoverImageSize))
	CoverSuffix         = ".cover"
	ProfileSuffix       = ".profile"
	ctx                 = context.Background()
)

func SetAccountProfileImage(accountAddress string, accountId uint64, image *string) (string, error) {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize == 0 {
		return "", errorImageZeroLen
	}
	if imageSize > maxProfileImageSize {
		return "", errorProfileTooBig
	}

	uploader, err := cdn.GetImageUploaderOrErr()
	if err != nil {
		return "", err
	}
	imgId := accountAddress + ProfileSuffix
	imgUrl, err := uploader.UploadBase64(ctx, (*image)[:len(*image)-1], imgId)
	if err != nil {
		return "", err
	}

	err = storage.UpdateAccountProfileWhereId(accountId, imgUrl)
	if err != nil {
		return "", err
	}

	return imgUrl, nil
}

func SetAccountCoverImage(accountAddress string, accountId uint64, image *string) (string, error) {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize == 0 {
		return "", errorImageZeroLen
	}
	if imageSize > maxCoverImageSize {
		return "", errorCoverTooBig
	}

	uploader, err := cdn.GetImageUploaderOrErr()
	if err != nil {
		return "", err
	}
	imgId := accountAddress + CoverSuffix
	imgUrl, err := uploader.UploadBase64(ctx, (*image)[:len(*image)-1], imgId)
	if err != nil {
		return "", err
	}

	err = storage.UpdateAccountCoverWhereId(accountId, imgUrl)
	if err != nil {
		return "", err
	}

	return imgUrl, nil
}

func SetCollectionCoverImage(tokenId string, collectionId uint64, image *string) (string, error) {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize == 0 {
		return "", errorImageZeroLen
	}
	if imageSize > maxCoverImageSize {
		return "", errorCoverTooBig
	}

	uploader, err := cdn.GetImageUploaderOrErr()
	if err != nil {
		return "", err
	}
	imgId := tokenId + CoverSuffix
	imgUrl, err := uploader.UploadBase64(ctx, (*image)[:len(*image)-1], imgId)
	if err != nil {
		return "", err
	}

	err = storage.UpdateCollectionCoverWhereId(collectionId, imgUrl)
	if err != nil {
		return "", err
	}

	return imgUrl, nil
}

func SetCollectionProfileImage(tokenId string, collectionId uint64, image *string) (string, error) {
	imageSize := getByteArrayLenOfBase64EncodedImage(image)
	if imageSize == 0 {
		return "", errorImageZeroLen
	}
	if imageSize > maxProfileImageSize {
		return "", errorCoverTooBig
	}

	uploader, err := cdn.GetImageUploaderOrErr()
	if err != nil {
		return "", err
	}
	imgId := tokenId + ProfileSuffix
	imgUrl, err := uploader.UploadBase64(ctx, (*image)[:len(*image)-1], imgId)
	if err != nil {
		return "", err
	}

	err = storage.UpdateCollectionProfileWhereId(collectionId, imgUrl)
	if err != nil {
		return "", err
	}

	return imgUrl, nil
}

func getByteArrayLenOfBase64EncodedImage(image *string) int {
	return base64.StdEncoding.DecodedLen(len(*image))
}
