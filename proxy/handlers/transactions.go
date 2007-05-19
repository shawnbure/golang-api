package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseTransactionsEndpoint         = "/transactions"
	transactionsListEndpoint         = "/list/:offset/:limit"
	transactionsByTokenEndpoint      = "/token/:tokenId/:offset/:limit"
	transactionsByAccountEndpoint    = "/account/:accountId/:offset/:limit"
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
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /transactions/token/{tokenId}/{offset}/{limit} [get]
func (handler *transactionsHandler) getByToken(c *gin.Context) {
	tokenIdString := c.Param("tokenId")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	tokenId, err := strconv.ParseUint(tokenIdString, 10, 64)
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

	transactions, err := storage.GetTransactionsByTokenIdWithOffsetLimit(tokenId, int(offset), int(limit))
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
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
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

	transactions, err := storage.GetTransactionsByBuyerOrSellerIdWithOffsetLimit(accountId, int(offset), int(limit))
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
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
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

	transactions, err := storage.GetTransactionsByCollectionIdWithOffsetLimit(collectionId, int(offset), int(limit))
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
