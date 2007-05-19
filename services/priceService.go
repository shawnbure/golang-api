package services

import (
	"encoding/json"
	"errors"
	"github.com/erdsea/erdsea-api/cache"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	EGLDPriceCacheKey           = "EGLDPrice"
	binanceEGLDPriceUrl         = "https://api.binance.com/api/v3/ticker/price?symbol=EGLDUSDT"
	binanceResponseExpirePeriod = 10 * time.Minute
)

type BinancePriceRequest struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func GetEGLDPrice() (float64, error) {
	var price float64

	errRead := cache.GetCacher().Get(EGLDPriceCacheKey, &price)
	if errRead == nil {
		return price, nil
	}

	var bpr BinancePriceRequest
	err := HttpGet(binanceEGLDPriceUrl, &bpr)
	if err != nil {
		log.Debug("binance request failed")
		return -1, err
	}
	if bpr.Price == "" {
		log.Debug("price is empty")
		return -1, errors.New("invalid response")
	}

	price, err = StrToFloat64(bpr.Price)
	if err != nil {
		log.Debug("could not parse price")
		return -1, errors.New("invalid response")
	}

	errSet := cache.GetCacher().Set(EGLDPriceCacheKey, price, binanceResponseExpirePeriod)
	if errSet != nil {
		log.Debug("could not cache result", errSet)
	}

	return price, err
}

func HttpGet(url string, castTarget interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return json.Unmarshal(respBytes, castTarget)
}

func StrToFloat64(v string) (float64, error) {
	vFloat, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return -1, err
	}

	return vFloat, nil
}
