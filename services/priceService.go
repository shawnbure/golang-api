package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/erdsea/erdsea-api/cache"
)

const (
	USD                         = "USD"
	USDTTicker                  = "USDT"
	EGLDTicker                  = "EGLD"
	EGLDPriceCacheKey           = "EGLDPrice"
	binancePriceUrl             = "https://api.binance.com/api/v3/ticker/price?symbol=%s%s"
	binanceResponseExpirePeriod = 10 * time.Minute
)

type BinancePriceRequest struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func GetPrice(from, to string) (float64, error) {
	to = strings.ToUpper(to)
	if strings.Contains(to, USD) {
		to = USDTTicker
	}

	from = strings.ToUpper(from)

	url := fmt.Sprintf(binancePriceUrl, from, to)

	var bpr BinancePriceRequest
	err := HttpGet(url, &bpr)
	if err != nil {
		log.Debug("binance request failed")
		return -1, err
	}
	if bpr.Price == "" {
		log.Debug("price is empty")
		return -1, errors.New("invalid response")
	}

	price, err := StrToFloat64(bpr.Price)
	if err != nil {
		log.Debug("could not parse price")
		return -1, errors.New("could not parse price")
	}

	return price, err
}

func GetEGLDPrice() (float64, error) {
	localCacher := cache.GetLocalCacher()
	priceVal, errRead := localCacher.Get(EGLDPriceCacheKey)
	if errRead == nil {
		return priceVal.(float64), nil
	}

	price, err := GetPrice(EGLDTicker, USDTTicker)
	if err != nil {
		return price, err
	}

	errSet := localCacher.SetWithTTLSync(EGLDPriceCacheKey, price, binanceResponseExpirePeriod)
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

func HttpGetRaw(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBytes), err
}

func StrToFloat64(v string) (float64, error) {
	vFloat, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return -1, err
	}

	return vFloat, nil
}
