package handlers

import (
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	baseStatsEndpoint                   = "/stats"
	StatTransactionsCountEndpoint       = "/txCount"
	StatTransactionsCountByDateEndpoint = "/txCount/:date"
)

type statsHandler struct {
}

func NewStatsHandler(groupHandler *groupHandler) {
	handler := &statsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: StatTransactionsCountEndpoint, HandlerFunc: handler.getTradeCounts},
		{Method: http.MethodGet, Path: StatTransactionsCountByDateEndpoint, HandlerFunc: handler.getTradeCounts},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseStatsEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Gets transactions count.
// @Description Gets transactions count (total/buy/withdraw/...) and can be filtered by date
// @Tags transactions
// @Accept json
// @Produce json
// @Param collectionId path string true "collection id"
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []dtos.TradesCount
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /stats/total/{date} [get]
func (handler *statsHandler) getTradeCounts(c *gin.Context) {
	filterDate := c.Param("date")

	result := dtos.TradesCount{}
	if strings.TrimSpace(filterDate) == "" {
		totalCount, err := storage.GetTransactionsCount()
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		result.Total = totalCount

		totalBuy, err := storage.GetTransactionsCountByType(entities.BuyToken)
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		result.Buy = totalBuy

		totalWithdraw, err := storage.GetTransactionsCountByType(entities.WithdrawToken)
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		result.Withdraw = totalWithdraw
	} else {
		totalCount, err := storage.GetTransactionsCountByDate(filterDate)
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		result.Total = totalCount

		totalBuy, err := storage.GetTransactionsCountByDateAndType(entities.BuyToken, filterDate)
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		result.Buy = totalBuy

		totalWithdraw, err := storage.GetTransactionsCountByDateAndType(entities.WithdrawToken, filterDate)
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		result.Withdraw = totalWithdraw

	}

	dtos.JsonResponse(c, http.StatusOK, result, "")
}
