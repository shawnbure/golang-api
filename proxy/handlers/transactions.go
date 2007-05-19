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
	baseTransactionsEndpoint         = "/transactions"
	getTransactionsEndpoint          = "/:offset/:limit"
	getTransactionsByAssetEndpoint   = "/by-asset/:tokenId/:nonce/:offset/:limit"
	getTransactionsByAddressEndpoint = "/by-address/:address/:offset/:limit"
)

type transactionsHandler struct {
}

func NewTransactionsHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &transactionsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: getTransactionsEndpoint, HandlerFunc: handler.get},
		{Method: http.MethodGet, Path: getTransactionsByAddressEndpoint, HandlerFunc: handler.getByAddress},
		{Method: http.MethodGet, Path: getTransactionsByAssetEndpoint, HandlerFunc: handler.getByAsset},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseTransactionsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *transactionsHandler) get(c *gin.Context) {
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

	transactions, err := storage.GetTransactionsWithOffsetLimit(offset, limit)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, transactions, "")
}

func (handler *transactionsHandler) getByAsset(c *gin.Context) {
	tokenId := c.Param("tokenId")
	nonceString := c.Param("nonce")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	nonce, err := strconv.ParseUint(nonceString, 10, 64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

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

	asset, err := storage.GetAssetByTokenIdAndNonce(tokenId, nonce)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	transactions, err := storage.GetTransactionsByAssetIdWithOffsetLimit(asset.ID, offset, limit)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, transactions, "")
}

func (handler *transactionsHandler) getByAddress(c *gin.Context) {
	address := c.Param("address")
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

	account, err := storage.GetAccountByAddress(address)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	transactions, err := storage.GetTransactionsByBuyerOrSellerIdWithOffsetLimit(account.ID, offset, limit)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, transactions, "")
}
