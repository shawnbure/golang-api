package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ENFT-DAO/youbei-api/data/entities"
)

const (
	itemsBaseUrl			=	"https://api.dreamship.com/v1/items"
	availableItemsUrl		=	"%s/%d/"
	shippingUrl				=	"%s/%d/%s/"
	usShippingUrl			=	"us-shipping-methods"
	internationalShipping	=	"international-shipping-methods"
)

// 19 is item id of canvas in dreamship
// To add more item, just add its id, can be find here https://api.dreamship.com/v1/items/
var availableItem = [1]int64{19}


func GetAvailableVariants() ([]entities.DreamshipItems, error) {
	var availableItems []entities.DreamshipItems
	for _, item := range availableItem{
		url := fmt.Sprintf(availableItemsUrl, itemsBaseUrl, item)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer 7a1aace437f7b50e420329bef9e7804f2cca65a7")
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

func GetShipmentMethodsAndCosts(contryCode string, stateCode string, item int64) (entities.ShippingMethodResponse, error) {
	
	url := fmt.Sprintf(shippingUrl, itemsBaseUrl, item, usShippingUrl)
	searachFor := stateCode
	if contryCode != "US"  {
		url = fmt.Sprintf(shippingUrl, itemsBaseUrl, item, internationalShipping)
		searachFor = contryCode
	}

	
	req, _ := http.NewRequest("GET", url, nil)
	
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer 7a1aace437f7b50e420329bef9e7804f2cca65a7")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	
	var m map[string]entities.ShippingMethodResponse
	err := json.Unmarshal(body, &m)
	if err != nil {
		return entities.ShippingMethodResponse{}, err
	}

	v, found := m[searachFor]
	if !found {
		return entities.ShippingMethodResponse{}, nil
	}

	return v, nil
}