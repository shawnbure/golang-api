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
	baseReportEndpoint                                   = "/reports"
	ReportSalesLast24HoursOverallEndpoint                = "/sales/daily/overall"
	ReportSalesLast24HoursTransactionsEndpoint           = "/sales/daily/transactions"
	ReportSalesLastWeekTopSellerOverallEndpoint          = "/sales/weekly/top/overall"
	ReportSalesLastWeekTopSellerTransactionsEndpoint     = "/sales/weekly/top/transactions"
	ReportListingLast24HoursVerifiedTransactionsEndpoint = "/listing/daily/transactions/verified"
	ReportListingLast24HoursAllTransactionsEndpoint      = "/listing/daily/transactions/all"
	ReportBuyLastWeekTopSellerOverallEndpoint            = "/buy/weekly/top/overall"
	ReportBuyLastWeekTopSellerTransactionsEndpoint       = "/buy/weekly/top/transactions"
)

const (
	TopListThreshold = 10
)

type reportHandler struct {
}

func NewReportHandler(groupHandler *groupHandler) {
	handler := &reportHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: ReportSalesLast24HoursOverallEndpoint, HandlerFunc: handler.getLast24HoursSalesOverall},
		{Method: http.MethodGet, Path: ReportSalesLast24HoursTransactionsEndpoint, HandlerFunc: handler.getLast24HoursSalesTransactions},
		{Method: http.MethodGet, Path: ReportSalesLastWeekTopSellerOverallEndpoint, HandlerFunc: handler.getLastWeekTopSellerOverall},
		{Method: http.MethodGet, Path: ReportSalesLastWeekTopSellerTransactionsEndpoint, HandlerFunc: handler.getLastWeekTopSellerTransactions},
		{Method: http.MethodGet, Path: ReportBuyLastWeekTopSellerOverallEndpoint, HandlerFunc: handler.getLastWeekTopBuyerOverall},
		{Method: http.MethodGet, Path: ReportBuyLastWeekTopSellerTransactionsEndpoint, HandlerFunc: handler.getLastWeekTopBuyerTransactions},
		{Method: http.MethodGet, Path: ReportListingLast24HoursVerifiedTransactionsEndpoint, HandlerFunc: handler.getLast24HoursVerifiedListingTransactions},
		{Method: http.MethodGet, Path: ReportListingLast24HoursAllTransactionsEndpoint, HandlerFunc: handler.getLast24HoursListingTransactions},
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

// @Summary Gets Last Week Top Sellers
// @Description Gets Last Week Top Sellers
// @Tags reports
// @Accept json
// @Param format query string false  "format of the output"
// @Produce application/csv
// @Success 200 {file} file
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /reports/sales/weekly/top/overall [get]
func (handler *reportHandler) getLastWeekTopSellerOverall(c *gin.Context) {
	q := c.Request.URL.Query().Get("format")
	if strings.TrimSpace(q) == "" {
		q = "csv"
	} else {
		q = strings.TrimSpace(q)
	}

	currentTime := time.Now().UTC()
	oneWeekBefore := currentTime.Add(-24 * 7 * time.Hour)

	currentTimeStr := fmt.Sprintf("%4d-%02d-%02d 00:00:00", currentTime.Year(), currentTime.Month(), currentTime.Day())
	oneWeekBeforeStr := fmt.Sprintf("%4d-%02d-%02d 00:00:00", oneWeekBefore.Year(), oneWeekBefore.Month(), oneWeekBefore.Day())

	records, err := storage.GetTopBestSellerLastWeek(TopListThreshold, oneWeekBeforeStr, currentTimeStr)
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
		result := [][]string{}
		result = append(result, []string{
			"From Time", "To Time", "Address", "Volume",
		})

		for _, record := range records {
			result = append(result, []string{
				oneWeekBeforeStr, currentTimeStr, record.Address, fmt.Sprintf("%f", record.Volume),
			})
		}

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
		for index, rec := range records {
			rec.FromTime = oneWeekBeforeStr
			rec.ToTime = currentTimeStr
			records[index] = rec
		}

		result := dtos.ReportTopVolumeByAddress{Records: records}
		dtos.JsonResponse(c, http.StatusOK, result, "")
	} else {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
	}
}

// @Summary Gets Last Week Top Sellers Transactions
// @Description Gets Last Week Top Sellers Transactions
// @Tags reports
// @Accept json
// @Param format query string false  "format of the output"
// @Produce application/csv
// @Success 200 {file} file
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /reports/sales/weekly/top/transactions [get]
func (handler *reportHandler) getLastWeekTopSellerTransactions(c *gin.Context) {
	q := c.Request.URL.Query().Get("format")
	if strings.TrimSpace(q) == "" {
		q = "csv"
	} else {
		q = strings.TrimSpace(q)
	}

	currentTime := time.Now().UTC()
	oneWeekBefore := currentTime.Add(-24 * 7 * time.Hour)

	currentTimeStr := fmt.Sprintf("%4d-%02d-%02d 00:00:00", currentTime.Year(), currentTime.Month(), currentTime.Day())
	oneWeekBeforeStr := fmt.Sprintf("%4d-%02d-%02d 00:00:00", oneWeekBefore.Year(), oneWeekBefore.Month(), oneWeekBefore.Day())

	records, err := storage.GetTopBestSellerLastWeek(TopListThreshold, oneWeekBeforeStr, currentTimeStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	addresses := []string{}
	for _, r := range records {
		addresses = append(addresses, r.Address)
	}

	transactions, err := storage.GetTopBestSellerLastWeekTransactions(oneWeekBeforeStr, currentTimeStr, addresses)
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
		result := dtos.ReportTopVolumeByAddressTransactionsList{Transactions: transactions}
		dtos.JsonResponse(c, http.StatusOK, result, "")
	} else {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
	}
}

// @Summary Gets Last Week Top Buyers
// @Description Gets Last Week Top Buyers
// @Tags reports
// @Accept json
// @Param format query string false  "format of the output"
// @Produce application/csv
// @Success 200 {file} file
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /reports/buy/weekly/top/overall [get]
func (handler *reportHandler) getLastWeekTopBuyerOverall(c *gin.Context) {
	q := c.Request.URL.Query().Get("format")
	if strings.TrimSpace(q) == "" {
		q = "csv"
	} else {
		q = strings.TrimSpace(q)
	}

	currentTime := time.Now().UTC()
	oneWeekBefore := currentTime.Add(-24 * 7 * time.Hour)

	currentTimeStr := fmt.Sprintf("%4d-%02d-%02d 00:00:00", currentTime.Year(), currentTime.Month(), currentTime.Day())
	oneWeekBeforeStr := fmt.Sprintf("%4d-%02d-%02d 00:00:00", oneWeekBefore.Year(), oneWeekBefore.Month(), oneWeekBefore.Day())

	records, err := storage.GetTopBestBuyerLastWeek(TopListThreshold, oneWeekBeforeStr, currentTimeStr)
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
		result := [][]string{}
		result = append(result, []string{
			"From Time", "To Time", "Address", "Volume",
		})

		for _, record := range records {
			result = append(result, []string{
				oneWeekBeforeStr, currentTimeStr, record.Address, fmt.Sprintf("%f", record.Volume),
			})
		}

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
		for index, rec := range records {
			rec.FromTime = oneWeekBeforeStr
			rec.ToTime = currentTimeStr
			records[index] = rec
		}

		result := dtos.ReportTopVolumeByAddress{Records: records}
		dtos.JsonResponse(c, http.StatusOK, result, "")
	} else {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
	}
}

// @Summary Gets Last Week Top Buyers Transactions
// @Description Gets Last Week Top Buyers Transactions
// @Tags reports
// @Accept json
// @Param format query string false  "format of the output"
// @Produce application/csv
// @Success 200 {file} file
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /reports/buy/weekly/top/transactions [get]
func (handler *reportHandler) getLastWeekTopBuyerTransactions(c *gin.Context) {
	q := c.Request.URL.Query().Get("format")
	if strings.TrimSpace(q) == "" {
		q = "csv"
	} else {
		q = strings.TrimSpace(q)
	}

	currentTime := time.Now().UTC()
	oneWeekBefore := currentTime.Add(-24 * 7 * time.Hour)

	currentTimeStr := fmt.Sprintf("%4d-%02d-%02d 00:00:00", currentTime.Year(), currentTime.Month(), currentTime.Day())
	oneWeekBeforeStr := fmt.Sprintf("%4d-%02d-%02d 00:00:00", oneWeekBefore.Year(), oneWeekBefore.Month(), oneWeekBefore.Day())

	records, err := storage.GetTopBestBuyerLastWeek(TopListThreshold, oneWeekBeforeStr, currentTimeStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	addresses := []string{}
	for _, r := range records {
		addresses = append(addresses, r.Address)
	}

	transactions, err := storage.GetTopBestBuyerLastWeekTransactions(oneWeekBeforeStr, currentTimeStr, addresses)
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
		result := dtos.ReportTopVolumeByAddressTransactionsList{Transactions: transactions}
		dtos.JsonResponse(c, http.StatusOK, result, "")
	} else {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
	}
}

// @Summary Gets Last 24 Hours Verified Listing Transactions
// @Description Gets Last 24 Hours Verified Listing Transactions
// @Tags reports
// @Accept json
// @Param format query string false  "format of the output"
// @Produce application/csv
// @Success 200 {file} file
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /listing/daily/transactions/verified [get]
func (handler *reportHandler) getLast24HoursVerifiedListingTransactions(c *gin.Context) {
	q := c.Request.URL.Query().Get("format")
	if strings.TrimSpace(q) == "" {
		q = "csv"
	} else {
		q = strings.TrimSpace(q)
	}

	currentTime := time.Now().UTC()
	oneDayBefore := currentTime.Add(-24 * time.Hour)

	//currentTimeStr := fmt.Sprintf("%4d-%02d-%02d %02d:00:00", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour())
	oneDayBeforeStr := fmt.Sprintf("%4d-%02d-%02d", oneDayBefore.Year(), oneDayBefore.Month(), oneDayBefore.Day())

	transactions, err := storage.GetListingTransactionsPerSpecificDay(oneDayBeforeStr, true)
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
			"Tx Id", "Tx Hash", "Tx Type", "Tx Price", "Tx Timestamp", "Address", "Token Id", "Token Name", "Token Image Link", "Collection Token Id", "Collection Name",
		})
		for _, item := range transactions {
			d := []string{
				fmt.Sprintf("%d", item.ID),
				item.Hash,
				string(item.Type),
				fmt.Sprintf("%f", item.PriceNominal),
				time.Unix(int64(item.Timestamp), 0).String(),
				item.Seller.Address,
				item.Token.TokenID,
				item.Token.TokenName,
				item.Token.ImageLink,
				item.Collection.CollectionTokenID,
				item.Collection.Name,
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
		result := dtos.ReportLast24HoursVerifiedListingTransactionsList{Transactions: transactions}
		dtos.JsonResponse(c, http.StatusOK, result, "")
	} else {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
	}
}

// @Summary Gets Last 24 Hours Listing Transactions
// @Description Gets Last 24 Hours Listing Transactions
// @Tags reports
// @Accept json
// @Param format query string false  "format of the output"
// @Produce application/csv
// @Success 200 {file} file
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /listing/daily/transactions/all [get]
func (handler *reportHandler) getLast24HoursListingTransactions(c *gin.Context) {
	q := c.Request.URL.Query().Get("format")
	if strings.TrimSpace(q) == "" {
		q = "csv"
	} else {
		q = strings.TrimSpace(q)
	}

	currentTime := time.Now().UTC()
	oneDayBefore := currentTime.Add(-24 * time.Hour)

	//currentTimeStr := fmt.Sprintf("%4d-%02d-%02d %02d:00:00", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour())
	oneDayBeforeStr := fmt.Sprintf("%4d-%02d-%02d", oneDayBefore.Year(), oneDayBefore.Month(), oneDayBefore.Day())

	transactions, err := storage.GetListingTransactionsPerSpecificDay(oneDayBeforeStr, false)
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
			"Tx Id", "Tx Hash", "Tx Type", "Tx Price", "Tx Timestamp", "Address", "Token Id", "Token Name", "Token Image Link", "Collection Token Id", "Collection Name", "Collection Verified",
		})
		for _, item := range transactions {
			d := []string{
				fmt.Sprintf("%d", item.ID),
				item.Hash,
				string(item.Type),
				fmt.Sprintf("%f", item.PriceNominal),
				time.Unix(int64(item.Timestamp), 0).String(),
				item.Seller.Address,
				item.Token.TokenID,
				item.Token.TokenName,
				item.Token.ImageLink,
				item.Collection.CollectionTokenID,
				item.Collection.Name,
				fmt.Sprintf("%v", item.Collection.IsVerified),
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
		result := dtos.ReportLast24HoursVerifiedListingTransactionsList{Transactions: transactions}
		dtos.JsonResponse(c, http.StatusOK, result, "")
	} else {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
	}
}
