package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/process"
	"github.com/gin-gonic/gin"
)

const (
	baseEventsEndpoint    = "/events"
	pushEventsEndpoint    = "/push"
	pushFinalizedEndpoint = "/finalized"
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
		{Method: http.MethodPost, Path: pushFinalizedEndpoint, HandlerFunc: h.pushFinalizedEvents},
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
	var blockEvents entities.BlockEvents

	err := c.Bind(&blockEvents)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	if blockEvents.Events != nil {
		h.processor.OnEvents(blockEvents)
	}

	dtos.JsonResponse(c, http.StatusOK, nil, "")
}

func (h *eventsHandler) pushFinalizedEvents(c *gin.Context) {
	var finalizedBlock entities.FinalizedBlock

	err := c.Bind(&finalizedBlock)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	log.Println("received finalized events at:", time.Now().Unix())
	log.Println("finalized hash:", finalizedBlock.Hash)

	dtos.JsonResponse(c, http.StatusOK, nil, "")
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
