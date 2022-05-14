package handlers

import (
	"net/http"

	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseDreamshipUrl	=	"/print"
	availableItemsUrl	=	"/available_items"
	itemsVariantsUrl	=	"/item_variants"
	shippingStatusUrl	=	"/shipping_status/:countryCode/:stateCode"
)

type dreamshipHandler struct {
}

func NewDreamshipHandler(groupHandler *groupHandler) {
	handler := &dreamshipHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: shippingStatusUrl, HandlerFunc: handler.getShippingStatus},
		{Method: http.MethodGet, Path: availableItemsUrl, HandlerFunc: handler.getAvailableItems},
	}

	endpointGroupHandler := EndpointGroupHandler {
		Root: baseDreamshipUrl,
		Middlewares: []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (d *dreamshipHandler) getAvailableItems(c *gin.Context) {
	data, err := services.GetAvailableVariants()
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, "Cannot Fetch Data")
		return
	}
	dtos.JsonResponse(c, http.StatusOK, data, "")
}

func (d *dreamshipHandler) getShippingStatus(c *gin.Context) {
	countryCode := c.Param("countryCode")
	stateCode := c.Param("stateCode")
	data, err := services.GetShipmentMethodsAndCosts(countryCode, stateCode, 19)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, "can not fetch data")
		return
	}
	if data.Code == "" {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, "Country Code or State Code doesn't exist")
		return
	}
	if len(data.Methods) == 0 {
		dtos.JsonResponse(c, http.StatusOK, "Unfortunately there is no shipping method for your location.", "")
		return
	}
	dtos.JsonResponse(c, http.StatusOK, data, "")
}

