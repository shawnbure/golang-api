package handlers

import (
	"fmt"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/ENFT-DAO/youbei-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	baseReportEndpoint                         = "/reports"
	ReportSalesLast24HoursOverallEndpoint      = "/sales/daily/overall"
	ReportSalesLast24HoursTransactionsEndpoint = "/sales/daily/transactions"
)

type reportHandler struct {
}

func NewReportHandler(groupHandler *groupHandler) {
	handler := &reportHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: ReportSalesLast24HoursOverallEndpoint, HandlerFunc: handler.getLast24HoursSalesOverall},
		{Method: http.MethodGet, Path: ReportSalesLast24HoursTransactionsEndpoint, HandlerFunc: handler.getLast24HoursSalesTransactions},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseReportEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}
	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Gets Last 24 Hours Total Volume
// @Description Gets Last 24 Hours Total Volume
// @Tags reports
// @Accept json
// @Param format query string false  "format of the output"
// @Produce application/csv
// @Success 200 {file} file
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /reports/sales/daily/overall [get]
func (handler *reportHandler) getLast24HoursSalesOverall(c *gin.Context) {
	q := c.Request.URL.Query().Get("format")
	if strings.TrimSpace(q) == "" {
		q = "csv"
	} else {
		q = strings.TrimSpace(q)
	}

	currentTime := time.Now().UTC()
	oneDayBefore := currentTime.Add(-24 * time.Hour)

	currentTimeStr := fmt.Sprintf("%4d-%02d-%02d %02d:00:00", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour())
	oneDayBeforeStr := fmt.Sprintf("%4d-%02d-%02d %02d:00:00", oneDayBefore.Year(), oneDayBefore.Month(), oneDayBefore.Day(), oneDayBefore.Hour())

	totalC, _ := storage.GetLast24HoursTotalVolume(oneDayBeforeStr, currentTimeStr)

	transactions, err := storage.GetLast24HoursSalesTransactions(oneDayBeforeStr, currentTimeStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	transactionsLength := len(transactions)

	if q == "csv" || q == "raw" {

		csvWrapper, err := utils.NewCsvWrapper()
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		defer csvWrapper.Close()

		// Create csv wrapper and return the result
		result := [][]string{}
		result = append(result, []string{
			"From Time", "To Time", "Total Volume", "Total Transactions",
		})
		result = append(result, []string{
			oneDayBeforeStr, currentTimeStr, totalC.String(), strconv.Itoa(transactionsLength),
		})

		err = csvWrapper.WriteBulkRecord(result)
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}

		if q == "csv" {
			buff := csvWrapper.GetBuffer()
			dtos.ContentAsFileResponse(c, "data.csv", buff)
		} else {
			finalResult := csvWrapper.GetData()
			dtos.StringResponse(c, finalResult)
		}
	} else if q == "json" {
		f, _ := totalC.Float64()
		result := dtos.ReportLast24HoursOverall{
			FromTime:          oneDayBeforeStr,
			ToTime:            currentTimeStr,
			TotalVolume:       f,
			TotalVolumeStr:    totalC.String(),
			TotalTransactions: transactionsLength,
		}
		dtos.JsonResponse(c, http.StatusOK, result, "")
	} else {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
	}
}

// @Summary Gets Last 24 Hours Transactions
// @Description Gets Last 24 Hours Transactions
// @Tags reports
// @Accept json
// @Param format query string false  "format of the output"
// @Produce application/csv
// @Success 200 {file} file
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /reports/sales/daily/transactions [get]
func (handler *reportHandler) getLast24HoursSalesTransactions(c *gin.Context) {
	q := c.Request.URL.Query().Get("format")
	if strings.TrimSpace(q) == "" {
		q = "csv"
	} else {
		q = strings.TrimSpace(q)
	}

	currentTime := time.Now().UTC()
	oneDayBefore := currentTime.Add(-24 * time.Hour)

	currentTimeStr := fmt.Sprintf("%4d-%02d-%02d %02d:00:00", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour())
	oneDayBeforeStr := fmt.Sprintf("%4d-%02d-%02d %02d:00:00", oneDayBefore.Year(), oneDayBefore.Month(), oneDayBefore.Day(), oneDayBefore.Hour())

	transactions, err := storage.GetLast24HoursSalesTransactions(oneDayBeforeStr, currentTimeStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	if q == "csv" || q == "raw" {
		csvWrapper, err := utils.NewCsvWrapper()
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		defer csvWrapper.Close()

		// Create csv wrapper and return the result
		csvWrapper2, err := utils.NewCsvWrapper()
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		defer csvWrapper2.Close()

		result := [][]string{}
		result = append(result, []string{
			"Tx Hash", "Seller Address", "Buyer Address", "Token Id", "Token Name", "Token Media Link", "Price", "Time",
		})
		for _, item := range transactions {
			d := []string{
				item.TxHash,
				item.FromAddress,
				item.ToAddress,
				item.TokenId,
				item.TokenName,
				item.TokenImageLink,
				fmt.Sprintf("%f", item.TxPriceNominal),
				time.Unix(item.TxTimestamp, 0).String(),
			}
			result = append(result, d)
		}

		err = csvWrapper2.WriteBulkRecord(result)
		if err != nil {
			dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}

		if q == "csv" {
			buff2 := csvWrapper2.GetBuffer()
			dtos.ContentAsFileResponse(c, "data.csv", buff2)
		} else {
			finalResult := csvWrapper2.GetData()
			dtos.StringResponse(c, finalResult)
		}
	} else if q == "json" {
		result := dtos.ReportLast24HoursTransactionsList{Transactions: transactions}
		dtos.JsonResponse(c, http.StatusOK, result, "")
	} else {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
	}
}
