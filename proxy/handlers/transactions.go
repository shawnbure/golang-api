package handlers

import (
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseTransactionsEndpoint      = "/transactions"
	transactionsListEndpoint      = "/list/:offset/:limit"
	transactionsByAssetEndpoint   = "/asset/:tokenId/:nonce/:offset/:limit"
	transactionsByAddressEndpoint = "/address/:address/:offset/:limit"
)

type transactionsHandler struct {
}

func NewTransactionsHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &transactionsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: transactionsListEndpoint, HandlerFunc: handler.getList},
		{Method: http.MethodGet, Path: transactionsByAssetEndpoint, HandlerFunc: handler.getByAddress},
		{Method: http.MethodGet, Path: transactionsByAddressEndpoint, HandlerFunc: handler.getByAsset},
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

// @Summary Gets transaction for an asset.
// @Description Retrieves transactions for an asset. Unordered.
// @Tags transactions
// @Accept json
// @Produce json
// @Param tokenId path string true "token id"
// @Param nonce path int true "nonce"
// @Param offset path int true "offset"
// @Param limit path int true "limit"
// @Success 200 {object} []data.Transaction
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /transactions/asset/{tokenId}/{nonce}/{offset}/{limit} [get]
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

// @Summary Gets transaction for an address.
// @Description Retrieves transactions for an address. Unordered.
// @Tags transactions
// @Accept json
// @Produce json
// @Param address path int true "address"
// @Param offset path int true "offset"
// @Param limit path int true "limit"
// @Success 200 {object} []data.Transaction
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /transactions/address/{address}/{offset}/{limit} [get]
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
