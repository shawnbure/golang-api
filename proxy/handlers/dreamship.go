package handlers

import (
	"net/http"
	"strconv"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseDreamshipUrl	=	"/print"
	availableItemsUrl	=	"/available_items"
	shippingStatusUrl	=	"/shipping_status/:us_or_inter/:item_id"
	orderUrl			=	"/order/:walletAddress"
	orderByUserUrl		=	"/order/:walletAddress/:orderId"
	orderHookUrl		=	"/order/hook"
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
		{Method: http.MethodPost, Path: orderHookUrl, HandlerFunc: handler.setOrderHook},
		{Method: http.MethodGet, Path: orderUrl, HandlerFunc: handler.getOrdersList},
		{Method: http.MethodGet, Path: orderByUserUrl, HandlerFunc: handler.GetOrderByUser},
	}

	endpointGroupHandler := EndpointGroupHandler {
		Root: baseDreamshipUrl,
		Middlewares: []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *dreamshipHandler) setOrderHook(c *gin.Context) {
	var request = entities.ItemWebhook{}
	err := c.BindJSON(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, "", err.Error())
	}
	err = services.DreamshipWebHook(request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, "", err.Error())
	}
}

func (handler *dreamshipHandler) GetOrderByUser(c *gin.Context) {
	walletAddress := c.Param("walletAddress")
	orderId := c.Param("orderId")
	data, err := storage.RetrievesAnOrders(walletAddress, orderId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}
	dtos.JsonResponse(c, http.StatusAccepted, data, "")
}

func (handler *dreamshipHandler) getOrdersList(c *gin.Context) {
	walletAddress := c.Param("walletAddress")
	data, err := storage.RetrievesUserOrders(walletAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}
	dtos.JsonResponse(c, http.StatusAccepted, data, "")
}

func (handler *dreamshipHandler) setOrder(c *gin.Context) {
	var request = entities.DreamshipOrderItems{}
	walletAddress := c.Param("walletAddress")
	err := c.BindJSON(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, "", err.Error())
	}
	
	// Service Layer Should be added here.
	data, err := services.SetOrderHandler(handler.cfg, request, walletAddress)
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

