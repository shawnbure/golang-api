package handlers

import (
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	baseExplorerEndpoint      = "/explorer"
	ExplorerTokenListEndpoint = "/all/:timestamp/:currentPage/:nextPage"

	ExplorerPageSize = 20
)

type explorerHandler struct {
}

func NewExplorerHandler(groupHandler *groupHandler) {
	handler := &explorerHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: ExplorerTokenListEndpoint, HandlerFunc: handler.getExplorerTokensWithPagination},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseExplorerEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Gets Explorer Tokens With Pagination And Filtering
// @Description Gets Explorer Tokens With Pagination And Filtering
// @Tags explorer
// @Accept json
// @Produce json
// @Param timestamp path int64 true "last timestamp"
// @Param currentPage path int64 true "the current page"
// @Param nextPage path int64 true "the current page"
// @Param limit query int64 true "page size limit"
// @Param filter query string false  "filter parameter"
// @Param sort query string false  "sort option parameter"
// @Success 200 {object} dtos.ExplorerTokenList
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /explorer/all/{timestamp}/{currentPage}/{nextPage} [get]
func (handler *explorerHandler) getExplorerTokensWithPagination(c *gin.Context) {
	result := dtos.ExplorerTokenList{}

	timeSt := c.Param("timestamp")
	currentPageStr := c.Param("currentPage")
	nextPageStr := c.Param("nextPage")

	isVerifiedStr := c.Request.URL.Query().Get("verified")
	limit := c.Request.URL.Query().Get("limit")
	filter := c.Request.URL.Query().Get("filter")
	attrFilter := c.Request.URL.Query().Get("attrs")

	querySQL, queryValues, _ := services.ConvertFilterToQuery("tokens", filter)
	sqlFilter := entities.QueryFilter{Query: querySQL, Values: queryValues}

	sortStr := c.Request.URL.Query().Get("sort")
	sortSQL, sortValues, _ := services.ConvertSortToQuery("tokens", sortStr)
	sortOptions := entities.SortOptions{Query: sortSQL, Values: sortValues}

	attributes, _ := services.ConvertAttributeFilterToQuery(attrFilter)

	var ts int64 = 0
	var limitInt int = ExplorerPageSize

	ts, err := strconv.ParseInt(timeSt, 10, 64)
	if err != nil {
		ts = 0
	}

	currentPage, err := strconv.Atoi(currentPageStr)
	if err != nil {
		currentPage = 0
	}
	nextPage, err := strconv.Atoi(nextPageStr)
	if err != nil {
		nextPage = 0
	}

	limitInt, err = strconv.Atoi(limit)
	if err != nil || limitInt == 0 {
		limitInt = ExplorerPageSize
	}

	isVerified, err := strconv.ParseBool(isVerifiedStr)
	if err != nil {
		isVerified = false
	}

	tokens, totalCount, minPrice, maxPrice, err := services.GetAllExplorerTokens(services.GetAllExplorerTokensArgs{
		LastTimestamp: ts,
		Limit:         limitInt,
		Filter:        &sqlFilter,
		CurrentPage:   currentPage,
		NextPage:      nextPage,
		SortOptions:   &sortOptions,
		IsVerified:    isVerified,
		Attributes:    attributes,
	})
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	result.Tokens = tokens
	result.TotalCount = totalCount
	result.MinPrice = minPrice
	result.MaxPrice = maxPrice

	dtos.JsonResponse(c, http.StatusOK, result, "")
}
