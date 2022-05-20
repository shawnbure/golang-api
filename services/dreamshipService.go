package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/entities"
)

const (
	itemsBaseUrl						=	"https://api.dreamship.com/v1/items"
	orderUrl							= 	"https://api.dreamship.com/v1/orders/"
	availableItemsUrl					=	"%s/%d/"
	shippingUrl							=	"%s/%d/%s/"
	usShippingUrl						=	"us-shipping-methods"
	internationalShipping				=	"international-shipping-methods"

	dreamshipItemsCacheKey				=	"dreamshipItems"
	dreamshipShippingBaseCacheKey		=	"dreamshipShipping-%s"
	dreamshipItemsInfoExpirePeriod		=	6 * time.Hour
)

// 19 is item id of canvas in dreamship
// To add more item, just add its id, can be find here https://api.dreamship.com/v1/items/
var availableItem = [1]int64{19}

func SetOrderHandler(cfg config.ExternalCredentialConfig, order entities.DreamshipOrderItems) (interface{}, error) {
	// Call setOrder
	response, err := SetOrder(cfg, order)
	// manage error
	// Save setOrder with ID to user.
	return response, err
}

func DreamshipWebHook(order entities.ItemWebhook) error{
	// Update user order
	return nil
}

func SetOrder(cfg config.ExternalCredentialConfig, order entities.DreamshipOrderItems) (interface{}, error) {
	orderJson, err := json.Marshal(order)
	var response interface{}
	if err != nil {
		return response, err
	}
	req, _ := http.NewRequest("POST", orderUrl, bytes.NewBuffer(orderJson))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	bearer := fmt.Sprintf("Bearer %s", cfg.DreamshipAPIKey)
	req.Header.Add("Authorization", bearer)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, nil
	}
	return response, nil
}

func GetAvailableVariantsHandler(cfg config.ExternalCredentialConfig) ([]entities.DreamshipItems, error) {
	localCacher := cache.GetLocalCacher()

	itemsVal, errRead := localCacher.Get(dreamshipItemsCacheKey)
	if errRead == nil {
		return itemsVal.([]entities.DreamshipItems), nil
	}

	items, err := GetAvailableVariants(cfg)
	if err != nil {
		return items, err
	}
	errSet := localCacher.SetWithTTLSync(dreamshipItemsCacheKey, items, dreamshipItemsInfoExpirePeriod)
	if errSet != nil {
		log.Debug("could not cache result", errSet)
	}

	return items, nil
}

func GetShipmentMethodsAndCostsHandler(cfg config.ExternalCredentialConfig, usOrInternational string, item int64) (map[string]entities.ShippingMethodResponse, error) {
	localCacher := cache.GetLocalCacher()
	dreamshipShippingCacheKey := fmt.Sprintf(dreamshipShippingBaseCacheKey, usOrInternational)
	itemsVal, errRead := localCacher.Get(dreamshipShippingCacheKey)
	if errRead == nil {
		return itemsVal.(map[string]entities.ShippingMethodResponse), nil
	}

	items, err := GetShipmentMethodsAndCosts(cfg, usOrInternational, item)
	if err != nil {
		return items, err
	}
	errSet := localCacher.SetWithTTLSync(dreamshipShippingCacheKey, items, dreamshipItemsInfoExpirePeriod)
	if errSet != nil {
		log.Debug("could not cache result", errSet)
	}

	return items, nil
}

func GetAvailableVariants(cfg config.ExternalCredentialConfig) ([]entities.DreamshipItems, error) {
	var availableItems []entities.DreamshipItems
	for _, item := range availableItem{
		url := fmt.Sprintf(availableItemsUrl, itemsBaseUrl, item)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "application/json")
		bearer := fmt.Sprintf("Bearer %s", cfg.DreamshipAPIKey)
		req.Header.Add("Authorization", bearer)
		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		var availItem entities.DreamshipItems
		err := json.Unmarshal(body, &availItem)
		if err != nil {
			fmt.Println(err)
			return []entities.DreamshipItems{}, err
		}
		availableItems = append(availableItems, availItem) 
	}
	return availableItems, nil
}

func GetShipmentMethodsAndCosts(cfg config.ExternalCredentialConfig, usOrInternational string, item int64) (map[string]entities.ShippingMethodResponse, error) {
	
	url := fmt.Sprintf(shippingUrl, itemsBaseUrl, item, usShippingUrl)
	if usOrInternational != "us"  {
		url = fmt.Sprintf(shippingUrl, itemsBaseUrl, item, internationalShipping)
	}

	
	req, _ := http.NewRequest("GET", url, nil)
	
	req.Header.Add("Accept", "application/json")
	bearer := fmt.Sprintf("Bearer %s", cfg.DreamshipAPIKey)
	req.Header.Add("Authorization", bearer)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	
	var m map[string]entities.ShippingMethodResponse
	err := json.Unmarshal(body, &m)
	if err != nil {
		return map[string]entities.ShippingMethodResponse{}, err
	}

	return m, nil
}