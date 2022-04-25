package handlers

import (
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	baseActivityEndpoint  = "/activities"
	ActivitiesAllEndpoint = "/all/:timestamp/:id"
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
// @Param id path int64 true "last fetched id"
// @Success 200 {object} dtos.ActivityLogsList
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /activities/all/{timestamp}/{id} [get]
func (handler *activityHandler) getActivityListWithPagination(c *gin.Context) {
	result := dtos.ActivityLogsList{}

	timeSt := c.Param("timestamp")
	id := c.Param("id")

	var ts int64 = 0
	var lastId int64 = 0

	ts, err := strconv.ParseInt(timeSt, 10, 64)
	if err != nil {
		ts = 0
	}

	lastId, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		lastId = 0
	}

	transactions, err := storage.GetAllActivitiesWithPagination(lastId, ts, PageSize)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	result.Activities = transactions

	dtos.JsonResponse(c, http.StatusOK, result, "")
}
