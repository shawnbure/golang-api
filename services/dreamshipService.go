package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	dreamshipBaseUrl		=	"https://api.dreamship.com/v1/"
	shippingUrl				=	"https://api.dreamship.com/v1/items/%d/%s/"
	usShippingUrl			=	"us-shipping-methods"
	internationalShipping	=	"international-shipping-methods"
)

type ShippingMethodResponse struct {
	Code	string				`json:"code"`
	Name	string				`json:"name"`
	Methods	[]ShippingMethod	`json:"methods"`
}

type ShippingMethod struct {
	Cost			float64	`json:"cost"`
	Method			string	`json:"method"`
	DeliveryDaysMax	uint64	`json:"delivery_days_max"`
	DeliveryDaysMin	uint64	`json:"delivery_days_min"`
}

func GetShipmentMethodsAndCosts(contryCode string, stateCode string, item int64) (ShippingMethodResponse, error) {
	
	url := fmt.Sprintf(shippingUrl, item, usShippingUrl)
	searachFor := stateCode
	if contryCode != "US"  {
		url = fmt.Sprintf(shippingUrl, item, internationalShipping)
		searachFor = contryCode
		fmt.Println(url)
		fmt.Println(searachFor)
	}

	
	req, _ := http.NewRequest("GET", url, nil)
	
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer 7a1aace437f7b50e420329bef9e7804f2cca65a7")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	
	var m map[string]ShippingMethodResponse
	err := json.Unmarshal(body, &m)
	if err != nil {
		return ShippingMethodResponse{}, err
	}

	v, found := m[searachFor]
	if !found {
		return ShippingMethodResponse{}, nil
	}

	return v, nil
}