package services

import (
	"errors"
	"fmt"
	"github.com/erdsea/erdsea-api/data/images"
	"github.com/erdsea/erdsea-api/storage"
)

var (
	maxProfileImageSize = 512 * 1024
	maxCoverImageSize   = 1024 * 1024
	errorProfileTooBig  = errors.New(fmt.Sprintf("profile image exceeded max size of %d", maxProfileImageSize))
	errorCoverTooBig    = errors.New(fmt.Sprintf("profile image exceeded max size of %d", maxProfileImageSize))
)

func GetAccountProfileImage(userAddress string) (*string, error) {
	account, err := storage.GetAccountByAddress(userAddress)
	if err != nil {
		return nil, err
	}

	image, err := storage.GetAccountProfileImageByUserId(account.ID)
	if err != nil {
		return nil, err
	}

	return &image.ImageBase64, nil
}

func SetAccountProfileImage(userAddress string, image *string) error {
	if len(*image) > maxProfileImageSize {
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
	return storage.AddNewAccountProfileImageByUserId(&profileImage)
}

func GetAccountCoverImage(userAddress string) (*string, error) {
	account, err := storage.GetAccountByAddress(userAddress)
	if err != nil {
		return nil, err
	}

	image, err := storage.GetAccountCoverImageByUserId(account.ID)
	if err != nil {
		return nil, err
	}

	return &image.ImageBase64, nil
}

func SetAccountCoverImage(userAddress string, image *string) error {
	if len(*image) > maxCoverImageSize {
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
	return storage.AddNewAccountCoverImageByUserId(&coverImage)
}
