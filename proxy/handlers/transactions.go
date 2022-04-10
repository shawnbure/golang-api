package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/ENFT-DAO/youbei-api/stats/collstats"

	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseTransactionsEndpoint         = "/transactions"
	transactionsListEndpoint         = "/list/:offset/:limit"
	transactionsByTokenEndpoint      = "/token/:tokenId/:nonce/:offset/:limit"
	transactionsByAccountEndpoint    = "/account/:userAddress/:offset/:limit"
	transactionsByCollectionEndpoint = "/collection/:collectionId/:offset/:limit"
)

const MaxQueryGetLimit = 50

type transactionsHandler struct {
}

func NewTransactionsHandler(groupHandler *groupHandler) {
	handler := &transactionsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: transactionsListEndpoint, HandlerFunc: handler.getList},
		{Method: http.MethodGet, Path: transactionsByTokenEndpoint, HandlerFunc: handler.getByToken},
		{Method: http.MethodGet, Path: transactionsByAccountEndpoint, HandlerFunc: handler.getByAccount},
		{Method: http.MethodGet, Path: transactionsByCollectionEndpoint, HandlerFunc: handler.getByCollection},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseTransactionsEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Gets transaction list.
// @Description Retrieves transactions. Unordered.
// @Tags transactions
// @Accept json
// @Produce json
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /transactions/list/{offset}/{limit} [get]
func (handler *transactionsHandler) getList(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	offset, err := strconv.ParseUint(offsetStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	err = ValidateLimit(limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	transactions, err := storage.GetTransactionsWithOffsetLimit(int(offset), int(limit))
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, transactions, "")
}

// @Summary Gets transaction for an token.
// @Description Retrieves transactions for an token. Unordered.
// @Tags transactions
// @Accept json
// @Produce json
// @Param tokenId path string true "token id"
// @Param nonce path int true "nonce"
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /transactions/token/{tokenId}/{nonce}/{offset}/{limit} [get]
func (handler *transactionsHandler) getByToken(c *gin.Context) {
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	offset, err := strconv.ParseUint(offsetStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	err = ValidateLimit(limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	cacheInfo, err := services.GetOrAddTokenCacheInfo(tokenId, nonce)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	transactions, err := storage.GetTransactionsByTokenIdWithOffsetLimit(cacheInfo.TokenDbId, int(offset), int(limit))
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
// @Param userAddress path string true "user wallet address"
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /transactions/account/{userAddress}/{offset}/{limit} [get]
func (handler *transactionsHandler) getByAccount(c *gin.Context) {
	userAddress := c.Param("userAddress")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	offset, err := strconv.ParseUint(offsetStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	err = ValidateLimit(limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	cacheInfo, err := services.GetOrAddAccountCacheInfo(userAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	transactions, err := storage.GetTransactionsByBuyerOrSellerIdWithOffsetLimit(cacheInfo.AccountId, int(offset), int(limit))
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
// @Param collectionId path string true "collection id"
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /transactions/collection/{collectionId}/{offset}/{limit} [get]
func (handler *transactionsHandler) getByCollection(c *gin.Context) {
	tokenId := c.Param("collectionId")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	offset, err := strconv.ParseUint(offsetStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	err = ValidateLimit(limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	cacheInfo, err := collstats.GetOrAddCollectionCacheInfo(tokenId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	transactions, err := storage.GetTransactionsByCollectionIdWithOffsetLimit(cacheInfo.CollectionId, int(offset), int(limit))
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, transactions, "")
}

// @Summary Gets transactions count.
// @Description Gets transactions count (total/buy/withdraw/...) and can be filtered by date
// @Tags transactions
// @Accept json
// @Produce json
// @Param collectionId path string true "collection id"
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.TransactionsCount
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /transactions/total/{date} [get]
func (handler *transactionsHandler) getTradeCounts(c *gin.Context) {
	filterDate := c.Param("date")

	if strings.TrimSpace(filterDate) == "" {

	} else {

	}
	transactions, err := storage.GetTransactionsCount()
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, transactions, "")
}

func ValidateLimit(limit uint64) error {
	if limit > MaxQueryGetLimit {
		return errors.New("limit too big")
	}

	return nil
}
