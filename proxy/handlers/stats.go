package handlers

import (
	"fmt"
	"github.com/ENFT-DAO/youbei-api/services"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/gin-gonic/gin"
)

var (
	StatsTotalVolumeKeyFormat           = "Stats:Volume:Total"
	StatsTotalVolumeLastUpdateKeyFormat = "Stats:Volume:TotalLastUpdate"
	StatsTotalVolumeExpirePeriod        = 2 * time.Hour

	StatsTotalVolumePerDayKeyFormat    = "Stats:Volume:%s"
	StatsTotalVolumePerDayExpirePeriod = time.Hour * 24 * 1

	StatsTokensTotalCountKeyFormat    = "Stats:Tokens:TotalCount"
	StatsTokensTotalCountExpirePeriod = 15 * time.Minute
)

var (
	logInstance = logger.GetOrCreate("stats-handler")
)

const (
	baseStatsEndpoint                   = "/stats"
	StatTransactionsCountEndpoint       = "/txCount"
	StatTransactionsCountByDateEndpoint = "/txCount/:date"
	StatTotalVolumeEndpoint             = "/volume/total"
	StatTotalVolumeLastWeekPerDay       = "/volume/lastWeek"
	StatTokensTotalCount                = "/tokens/totalCount"
	StatListTransactionsWithPagination  = "/transactions/list/:timestamp/:limit"
)

const (
	PageSize = 20
)

type statsHandler struct {
}

func NewStatsHandler(groupHandler *groupHandler) {
	handler := &statsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: StatTransactionsCountEndpoint, HandlerFunc: handler.getTradeCounts},
		{Method: http.MethodGet, Path: StatTransactionsCountByDateEndpoint, HandlerFunc: handler.getTradeCounts},
		{Method: http.MethodGet, Path: StatTotalVolumeEndpoint, HandlerFunc: handler.getTotalTradesVolume},
		{Method: http.MethodGet, Path: StatTotalVolumeLastWeekPerDay, HandlerFunc: handler.getTotalTradesVolumeLastWeek},
		{Method: http.MethodGet, Path: StatTokensTotalCount, HandlerFunc: handler.getTokensTotalCount},
		{Method: http.MethodGet, Path: StatListTransactionsWithPagination, HandlerFunc: handler.getTransactionsListWithPagination},
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
// @Tags stats
// @Accept json
// @Produce json
// @Param date path string true "specific date"
// @Success 200 {object} dtos.TradesCount
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /stats/txCount/{date} [get]
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

		totalBuy, err := storage.GetTransactionsCountByType((entities.ListToken))
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		result.Buy = totalBuy

		totalWithdraw, err := storage.GetTransactionsCountByType((entities.WithdrawToken))
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

// @Summary Gets Total Volume
// @Description Gets Total Volume
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} dtos.TradesVolumeTotal
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /stats/volume/total [get]
func (handler *statsHandler) getTotalTradesVolume(c *gin.Context) {
	result := dtos.TradesVolumeTotal{}

	// Let's check the cache first
	localCacher := cache.GetLocalCacher()

	var totalVolume *big.Float
	var totalVolumeLastUpdate int64

	totalLU, errRead := localCacher.Get(StatsTotalVolumeLastUpdateKeyFormat)
	totalStr, errRead2 := localCacher.Get(StatsTotalVolumeKeyFormat)
	if errRead == nil && errRead2 == nil {
		totalVolume, _ = new(big.Float).SetString(totalStr.(string))
		totalVolumeLastUpdate = totalLU.(int64)
	} else {
		// get it from database and also cache it
		totalV, err := storage.GetTotalTradedVolume()
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		totalVolume = totalV
		totalVolumeLastUpdate = time.Now().UTC().Unix()

		err = localCacher.SetWithTTLSync(StatsTotalVolumeKeyFormat, totalV.String(), StatsTotalVolumeExpirePeriod)
		if err != nil {
			logInstance.Debug("could not set cache", "err", err)
		}

		err = localCacher.SetWithTTLSync(StatsTotalVolumeLastUpdateKeyFormat, totalVolumeLastUpdate, StatsTotalVolumeExpirePeriod)
		if err != nil {
			logInstance.Debug("could not set cache", "err", err)
		}
	}

	result.Sum = totalVolume.String()
	result.LastUpdate = totalVolumeLastUpdate

	dtos.JsonResponse(c, http.StatusOK, result, "")
}

// @Summary Gets Total Volume Per Day For Last Week
// @Description Gets Total Volume Per Day For Last Week
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} []dtos.TradesVolume
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /stats/volume/last_week [get]
func (handler *statsHandler) getTotalTradesVolumeLastWeek(c *gin.Context) {
	result := []dtos.TradesVolume{}

	// Let's find out today
	today := time.Now().UTC()
	//dateFormat := "2021-02-01"

	// Let's check the cache first
	localCacher := cache.GetLocalCacher()

	for i := 0; i < 7; i++ {
		tempDate := today.Add(-24 * time.Duration(i) * time.Hour)
		finalDate := fmt.Sprintf("%4d-%02d-%02d", tempDate.Year(), tempDate.Month(), tempDate.Day())

		var totalVolume *big.Float

		key := fmt.Sprintf(StatsTotalVolumePerDayKeyFormat, finalDate)
		totalStr, errRead := localCacher.Get(key)
		if errRead == nil {
			totalVolume, _ = new(big.Float).SetString(totalStr.(string))
		} else {
			// get it from database and also cache it
			totalV, err := storage.GetTotalTradedVolumeByDate(finalDate)
			if err != nil {
				totalVolume = big.NewFloat(0)
				//dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
				//return
			}
			totalVolume = totalV

			err = localCacher.SetWithTTLSync(StatsTotalVolumeKeyFormat, totalV.String(), StatsTotalVolumePerDayExpirePeriod)
			if err != nil {
				logInstance.Debug("could not set cache", "err", err)
			}
		}

		result = append(result, dtos.TradesVolume{
			Sum:  totalVolume.String(),
			Date: finalDate,
		})
	}

	dtos.JsonResponse(c, http.StatusOK, result, "")
}

// @Summary Gets Total Tokens Count
// @Description Gets Total Tokens Count
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} dtos.TokensTotalCount
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /stats/tokens/totalCount [get]
func (handler *statsHandler) getTokensTotalCount(c *gin.Context) {
	result := dtos.TokensTotalCount{}

	// Let's check the cache first
	localCacher := cache.GetLocalCacher()

	var totalCount int64

	totalC, errRead := localCacher.Get(StatsTokensTotalCountKeyFormat)
	if errRead == nil {
		totalCount = totalC.(int64)
	} else {
		// get it from database and also cache it
		totalC, err := storage.GetTotalTokenCount()
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		totalCount = totalC

		err = localCacher.SetWithTTLSync(StatsTokensTotalCountKeyFormat, totalCount, StatsTokensTotalCountExpirePeriod)
		if err != nil {
			logInstance.Debug("could not set cache", "err", err)
		}
	}

	result.Sum = totalCount
	dtos.JsonResponse(c, http.StatusOK, result, "")
}

// @Summary Gets Transactions List With Pagination
// @Description Gets Transactions List With Pagination
// @Tags stats
// @Accept json
// @Produce json
// @Param timestamp path int64 true "last timestamp"
// @Param limit path int64 true "page size limit"
// @Param filter query string false  "filter parameter"
// @Success 200 {object} dtos.StatTransactionsList
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /stats/transactions/list/{timestamp}/{limit} [get]
func (handler *statsHandler) getTransactionsListWithPagination(c *gin.Context) {
	result := dtos.StatTransactionsList{}

	timeSt := c.Param("timestamp")
	limit := c.Param("limit")

	filter := c.Request.URL.Query().Get("filter")
	querySQL, queryValues, _ := services.ConvertFilterToQuery("collections", filter)
	sqlFilter := entities.QueryFilter{Query: querySQL, Values: queryValues}

	var ts int64 = 0
	var limitInt int = PageSize

	ts, err := strconv.ParseInt(timeSt, 10, 64)
	if err != nil {
		ts = 0
	}

	limitInt, err = strconv.Atoi(limit)
	if err != nil || limitInt == 0 {
		limitInt = PageSize
	}

	transactions, err := services.GetAllTransactionsWithPagination(services.GetAllTransactionsWithPaginationArgs{
		LastTimestamp: ts,
		Limit:         limitInt,
		Filter:        &sqlFilter,
	})
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	result.Transactions = transactions

	dtos.JsonResponse(c, http.StatusOK, result, "")
}
