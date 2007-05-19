package handlers

import (
	"net/http"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/process"
	"github.com/gin-gonic/gin"
)

const (
	baseEventsEndpoint = "/events"
	pushEventsEndpoint = "/push"
)

type eventsHandler struct {
	config    config.ConnectorApiConfig
	processor *process.EventProcessor
}

func NewEventsHandler(
	groupHandler *groupHandler,
	processor *process.EventProcessor,
	config config.ConnectorApiConfig,
) error {
	h := &eventsHandler{
		config:    config,
		processor: processor,
	}

	endpoints := []EndpointHandler{
		{Method: http.MethodPost, Path: pushEventsEndpoint, HandlerFunc: h.pushEvents},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseEventsEndpoint,
		Middlewares:      h.createMiddlewares(),
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)

	return nil
}

func (h *eventsHandler) pushEvents(c *gin.Context) {
	var events []data.Event

	err := c.Bind(&events)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	if events != nil {
		h.processor.OnEvents(events)
	}

	data.JsonResponse(c, http.StatusOK, nil, "")
}

func (h *eventsHandler) createMiddlewares() []gin.HandlerFunc {
	var middleware []gin.HandlerFunc

	if h.config.Username != "" && h.config.Password != "" {
		basicAuth := gin.BasicAuth(gin.Accounts{
			h.config.Username: h.config.Password,
		})
		middleware = append(middleware, basicAuth)
	}

	return middleware
}
