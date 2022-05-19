package handlers

import (
	"net/http"
	"strconv"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseDreamshipUrl	=	"/print"
	availableItemsUrl	=	"/available_items"
	shippingStatusUrl	=	"/shipping_status/:us_or_inter/:item_id"
	orderUrl			=	"/order"
)

type dreamshipHandler struct {
	cfg	config.ExternalCredentialConfig
}

func NewDreamshipHandler(groupHandler *groupHandler, cfg config.ExternalCredentialConfig) {
	handler := &dreamshipHandler{cfg: cfg}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: shippingStatusUrl, HandlerFunc: handler.getShippingStatus},
		{Method: http.MethodGet, Path: availableItemsUrl, HandlerFunc: handler.getAvailableItems},
		{Method: http.MethodPost, Path: orderUrl, HandlerFunc: handler.setOrder},
	}

	endpointGroupHandler := EndpointGroupHandler {
		Root: baseDreamshipUrl,
		Middlewares: []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *dreamshipHandler) setOrder(c *gin.Context) {
	var request = entities.DreamshipOrderItems{}
	err := c.BindJSON(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, "", err.Error())
	}
	
	// Service Layer Should be added here.
	data, err := services.SetOrder(handler.cfg, request)
	if err != nil{
		dtos.JsonResponse(c, http.StatusBadRequest, "", err.Error())
	}

	dtos.JsonResponse(c, http.StatusCreated, data, "")
}


func (handler *dreamshipHandler) getAvailableItems(c *gin.Context) {
	data, err := services.GetAvailableVariantsHandler(handler.cfg)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, "Cannot Fetch Data")
		return
	}
	dtos.JsonResponse(c, http.StatusOK, data, "")
}

func (handler *dreamshipHandler) getShippingStatus(c *gin.Context) {
	usOrInternational := c.Param("us_or_inter")
	itemId := c.Param("item_id")
	item, err := strconv.ParseInt(itemId, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, "Please provide correct id for item")
	}
	data, err := services.GetShipmentMethodsAndCostsHandler(handler.cfg, usOrInternational, item)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, "can not fetch data")
		return
	}
	dtos.JsonResponse(c, http.StatusOK, data, "")
}

