package handlers

import (
	"github.com/erdsea/erdsea-api/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	baseCollectionsEndpoint = "/collections"
	getCollectionsEndpoint = "/:offset/:limit"
	getCollectionByNameEndpoint  = "/by-name/:name"
)

type collectionsHandler struct {
}

func NewCollectionsHandler(groupHandler *groupHandler) {
	handler := &collectionsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: getCollectionsEndpoint, HandlerFunc: handler.get},
		{Method: http.MethodGet, Path: getCollectionByNameEndpoint, HandlerFunc: handler.getByName},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseCollectionsEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *collectionsHandler) get(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	collections, err := storage.GetCollectionsWithOffsetLimit(offset, limit)
	if err != nil {
		JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	JsonResponse(c, http.StatusOK, collections, "")
}

func (handler *collectionsHandler) getByName(c *gin.Context) {
	collectionName := c.Param("collectionName")

	asset, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	JsonResponse(c, http.StatusOK, asset, "")
}
