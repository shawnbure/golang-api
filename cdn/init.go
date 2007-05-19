package cdn

import (
	"errors"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"sync"

	"github.com/cloudinary/cloudinary-go"
	"github.com/erdsea/erdsea-api/config"
)

var (
	once        sync.Once
	cloudy      *cloudinary.Cloudinary
	imgUploader ImageUploader

	log = logger.GetOrCreate("cdn")
)

const (
	local     = "local"
	cloudyCDN = "cloudy"
)

func MakeCloudyCDN(cfg config.CDNConfig) {
	once.Do(func() {
		newCloudy, err := cloudinary.NewFromParams(
			cfg.Name,
			cfg.ApiKey,
			cfg.ApiSecret,
		)
		if err != nil {
			panic(err)
		}

		cloudy = newCloudy
	})
}

func GetCloudyCDNOrErr() (*cloudinary.Cloudinary, error) {
	if cloudy == nil {
		return nil, errors.New("cloudy cdn is not initialized")
	}

	return cloudy, nil
}

func makeUploader(cfg config.CDNConfig) (ImageUploader, error) {
	switch cfg.Selector {
	case local:
		return NewLocalUploader(cfg), nil
	case cloudyCDN:
		return NewCloudyUploader(cfg)
	default:
		return nil, errors.New("unknown selector provided")
	}
}
