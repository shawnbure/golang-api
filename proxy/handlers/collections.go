package handlers

import (
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/services"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseCollectionsEndpoint         = "/collections"
	getCollectionsEndpoint          = "/:offset/:limit"
	getCollectionByNameEndpoint     = "/by-name/:collectionName"
	createCollectionEndpoint        = "/create"
	getCollectionStatisticsEndpoint = "/statistics/:collectionName"
)

type collectionsHandler struct {
	blockchainCfg config.BlockchainConfig
}

func NewCollectionsHandler(groupHandler *groupHandler, authCfg config.AuthConfig, blockchainCfg config.BlockchainConfig) {
	handler := &collectionsHandler{blockchainCfg: blockchainCfg}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: getCollectionsEndpoint, HandlerFunc: handler.get},
		{Method: http.MethodGet, Path: getCollectionByNameEndpoint, HandlerFunc: handler.getByName},
		{Method: http.MethodGet, Path: getCollectionStatisticsEndpoint, HandlerFunc: handler.getStatisticsForCollectionWithName},
		{Method: http.MethodPost, Path: createCollectionEndpoint, HandlerFunc: handler.create},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseCollectionsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (ch *collectionsHandler) get(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	collections, err := storage.GetCollectionsWithOffsetLimit(offset, limit)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, collections, "")
}

func (ch *collectionsHandler) getByName(c *gin.Context) {
	collectionName := c.Param("collectionName")

	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, collection, "")
}

func (ch *collectionsHandler) getStatisticsForCollectionWithName(c *gin.Context) {
	collectionName := c.Param("collectionName")

	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	stats, err := services.GetStatisticsForCollection(collection.ID)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, stats, "")
}

func (ch *collectionsHandler) create(c *gin.Context) {
	var request services.CreateCollectionRequest

	err := c.Bind(&request)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress, exists := c.Get(middleware.AddressKey)
	if !exists {
		data.JsonResponse(c, http.StatusInternalServerError, nil, "could not get address from context")
		return
	}
	if jwtAddress != request.UserAddress {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "jwt and request addresses differ")
		return
	}

	err = services.CreateCollection(&request, ch.blockchainCfg.ProxyUrl)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, "", "")
}
