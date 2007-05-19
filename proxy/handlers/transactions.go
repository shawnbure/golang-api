package handlers

import (
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseTransactionsEndpoint         = "/transactions"
	transactionsListEndpoint         = "/list/:offset/:limit"
	transactionsByAssetEndpoint      = "/asset/:assetId/:offset/:limit"
	transactionsByAccountEndpoint    = "/account/:accountId/:offset/:limit"
	transactionsByCollectionEndpoint = "/collection/:collectionId/:offset/:limit"
)

type transactionsHandler struct {
}

func NewTransactionsHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &transactionsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: transactionsListEndpoint, HandlerFunc: handler.getList},
		{Method: http.MethodGet, Path: transactionsByAssetEndpoint, HandlerFunc: handler.getByAsset},
		{Method: http.MethodGet, Path: transactionsByAccountEndpoint, HandlerFunc: handler.getByAccount},
		{Method: http.MethodGet, Path: transactionsByCollectionEndpoint, HandlerFunc: handler.getByCollection},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseTransactionsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Gets transaction list.
// @Description Retrieves transactions. Unordered.
// @Tags transactions
// @Accept json
// @Produce json
// @Param offset path int true "offset"
// @Param limit path int true "limit"
// @Success 200 {object} []data.Transaction
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /transactions/list/{offset}/{limit} [get]
func (handler *transactionsHandler) getList(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	transactions, err := storage.GetTransactionsWithOffsetLimit(offset, limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, transactions, "")
}

// @Summary Gets transaction for an asset.
// @Description Retrieves transactions for an asset. Unordered.
// @Tags transactions
// @Accept json
// @Produce json
// @Param assetId path uint64 true "asset id"
// @Param offset path int true "offset"
// @Param limit path int true "limit"
// @Success 200 {object} []data.Transaction
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /transactions/asset/{assetId}/{offset}/{limit} [get]
func (handler *transactionsHandler) getByAsset(c *gin.Context) {
	assetIdString := c.Param("assetId")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	assetId, err := strconv.ParseUint(assetIdString, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	transactions, err := storage.GetTransactionsByAssetIdWithOffsetLimit(assetId, offset, limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, transactions, "")
}

// @Summary Gets transaction for an account.
// @Description Retrieves transactions for an account. Unordered.
// @Tags transactions
// @Accept json
// @Produce json
// @Param accountId path uint64 true "account id"
// @Param offset path int true "offset"
// @Param limit path int true "limit"
// @Success 200 {object} []data.Transaction
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /transactions/account/{accountId}/{offset}/{limit} [get]
func (handler *transactionsHandler) getByAccount(c *gin.Context) {
	accountIdString := c.Param("accountId")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	accountId, err := strconv.ParseUint(accountIdString, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	transactions, err := storage.GetTransactionsByBuyerOrSellerIdWithOffsetLimit(accountId, offset, limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, transactions, "")
}

// @Summary Gets transaction for a collection.
// @Description Retrieves transactions for a collection. Unordered.
// @Tags transactions
// @Accept json
// @Produce json
// @Param collectionId path uint64 true "collection id"
// @Param offset path int true "offset"
// @Param limit path int true "limit"
// @Success 200 {object} []data.Transaction
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /transactions/collection/{collectionId}/{offset}/{limit} [get]
func (handler *transactionsHandler) getByCollection(c *gin.Context) {
	collectionIdString := c.Param("collectionId")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	collectionId, err := strconv.ParseUint(collectionIdString, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	transactions, err := storage.GetTransactionsByCollectionIdWithOffsetLimit(collectionId, offset, limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, transactions, "")
}
