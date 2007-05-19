package handlers

import (
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	baseCollectionsEndpoint     = "/collections"
	getCollectionsEndpoint      = "/:offset/:limit"
	getCollectionByNameEndpoint = "/by-name/:collectionName"
)

type collectionsHandler struct {
}

func NewCollectionsHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &collectionsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: getCollectionsEndpoint, HandlerFunc: handler.get},
		{Method: http.MethodGet, Path: getCollectionByNameEndpoint, HandlerFunc: handler.getByName},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseCollectionsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *collectionsHandler) get(c *gin.Context) {
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

func (handler *collectionsHandler) getByName(c *gin.Context) {
	collectionName := c.Param("collectionName")

	asset, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, asset, "")
}
