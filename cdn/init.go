package cdn

import (
	"errors"
	"sync"

	"github.com/cloudinary/cloudinary-go"
	"github.com/erdsea/erdsea-api/config"
)

var (
	once   sync.Once
	cloudy *cloudinary.Cloudinary
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
