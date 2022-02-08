package handlers

import (
	"net/http"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/proxy/middleware"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseSessionStatesEndpoint   = "/sessionStates"
	SessionStatesCreateEndpoint = "/create"
)

type stateSessionsHandler struct {
	blockchainCfg config.BlockchainConfig
}

func NewSessionStatesHandler(groupHandler *groupHandler, authCfg config.AuthConfig, blockchainCfg config.BlockchainConfig) {
	handler := &stateSessionsHandler{blockchainCfg: blockchainCfg}

	endpoints := []EndpointHandler{
		{Method: http.MethodPost, Path: SessionStatesCreateEndpoint, HandlerFunc: handler.create},
	}
	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseCollectionsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}
	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *stateSessionsHandler) create(c *gin.Context) {
	var request services.CreateSessionStateRequest

	collection, err := services.CreateSessionState(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, collection, "")
}
