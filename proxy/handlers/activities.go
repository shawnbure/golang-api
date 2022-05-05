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
	baseActivityEndpoint  = "/activities"
	ActivitiesAllEndpoint = "/all/:timestamp/:currentPage/:nextPage"

	ActivityPageSize = 7
)

type activityHandler struct {
}

func NewActivitiesHandler(groupHandler *groupHandler) {
	handler := &activityHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: ActivitiesAllEndpoint, HandlerFunc: handler.getActivityListWithPagination},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseActivityEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Gets Transactions Logs With Pagination
// @Description Gets Transactions Logs With Pagination
// @Tags activity
// @Accept json
// @Produce json
// @Param timestamp path int64 true "last timestamp"
// @Param currentPage path int64 true "the current page"
// @Param nextPage path int64 true "the current page"
// @Param timestamp path int64 true "last timestamp"
// @Param limit query int64 true "page size limit"
// @Param filter query string false  "filter parameter"
// @Success 200 {object} dtos.ActivityLogsList
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /activities/all/{timestamp}/{currentPage}/{nextPage} [get]
func (handler *activityHandler) getActivityListWithPagination(c *gin.Context) {
	result := dtos.ActivityLogsList{}

	timeSt := c.Param("timestamp")
	limit := c.Request.URL.Query().Get("limit")

	currentPageStr := c.Param("currentPage")
	nextPageStr := c.Param("nextPage")

	filter := c.Request.URL.Query().Get("filter")
	querySQL, queryValues, _ := services.ConvertFilterToQuery("transactions", filter)
	sqlFilter := entities.QueryFilter{Query: querySQL, Values: queryValues}

	var ts int64 = 0
	var limitInt int = ActivityPageSize

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
		limitInt = ActivityPageSize
	}

	transactions, err := services.GetAllActivities(services.GetAllActivityArgs{
		LastTimestamp: ts,
		Limit:         limitInt,
		Filter:        &sqlFilter,
		CurrentPage:   currentPage,
		NextPage:      nextPage,
	})
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	result.Activities = transactions
	dtos.JsonResponse(c, http.StatusOK, result, "")
}
