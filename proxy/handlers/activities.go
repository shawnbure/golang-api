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
	ActivitiesAllEndpoint = "/all/:timestamp/:limit"

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
// @Param limit path int64 true "page size limit"
// @Param filter query string false  "filter parameter"
// @Success 200 {object} dtos.ActivityLogsList
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /activities/all/{timestamp}/{limit} [get]
func (handler *activityHandler) getActivityListWithPagination(c *gin.Context) {
	result := dtos.ActivityLogsList{}

	timeSt := c.Param("timestamp")
	limit := c.Param("limit")

	filter := c.Request.URL.Query().Get("filter")
	querySQL, queryValues, _ := services.ConvertFilterToQuery("transactions", filter)
	sqlFilter := entities.QueryFilter{Query: querySQL, Values: queryValues}

	var ts int64 = 0
	var limitInt int = ActivityPageSize

	ts, err := strconv.ParseInt(timeSt, 10, 64)
	if err != nil {
		ts = 0
	}

	limitInt, err = strconv.Atoi(limit)
	if err != nil || limitInt == 0 {
		limitInt = ActivityPageSize
	}

	transactions, err := services.GetAllActivities(services.GetAllActivityArgs{
		LastTimestamp: ts,
		Limit:         limitInt,
		Filter:        &sqlFilter,
	})
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	result.Activities = transactions
	dtos.JsonResponse(c, http.StatusOK, result, "")
}
