package handlers

import (
	"github.com/erdsea/erdsea-api/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	baseAssetsEndpoint                = "/assets"
	getAssetByTokenIdAndNonceEndpoint = "/by-token-id-nonce/:tokenId/:nonce"
	getAssetsByCollectionEndpoint     = "/by-collection/:collectionName/:offset/:limit"
)

type assetsHandler struct {
}

func NewAssetsHandler(groupHandler *groupHandler) {
	handler := &assetsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: getAssetByTokenIdAndNonceEndpoint, HandlerFunc: handler.getByTokenIdAndNonce},
		{Method: http.MethodGet, Path: getAssetsByCollectionEndpoint, HandlerFunc: handler.getByCollection},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseAssetsEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *assetsHandler) getByTokenIdAndNonce(c *gin.Context) {
	tokenId := c.Param("tokenId")
	nonceString := c.Param("nonce")

	nonce, err := strconv.ParseUint(nonceString, 10, 64)
	if err != nil {
		JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	asset, err := storage.GetAssetByTokenIdAndNonce(tokenId, nonce)
	if err != nil {
		JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	JsonResponse(c, http.StatusOK, asset, "")
}

func (handler *assetsHandler) getByCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
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

	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	assets, err := storage.GetAssetsByCollectionIdWithOffsetLimit(collection.ID, offset, limit)
	if err != nil {
		JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	JsonResponse(c, http.StatusOK, assets, "")
}
