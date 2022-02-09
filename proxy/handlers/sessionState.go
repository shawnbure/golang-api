package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/proxy/middleware"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseSessionStatesEndpoint   = "/sessionStates"
	sessionStatesCreateEndpoint = "/create"
	sessionStatesDeleteEndpoint = "/delete/:accountId/:stateType"
)

type stateSessionsHandler struct {
	blockchainCfg config.BlockchainConfig
}

func NewSessionStatesHandler(groupHandler *groupHandler, authCfg config.AuthConfig, blockchainCfg config.BlockchainConfig) {
	handler := &stateSessionsHandler{blockchainCfg: blockchainCfg}

	endpoints := []EndpointHandler{
		{Method: http.MethodPost, Path: sessionStatesCreateEndpoint, HandlerFunc: handler.create},
		{Method: http.MethodPost, Path: sessionStatesDeleteEndpoint, HandlerFunc: handler.delete},
	}
	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseSessionStatesEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}
	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *stateSessionsHandler) create(c *gin.Context) {
	var request services.CreateSessionStateRequest

	fmt.Println("request.JsonData: " + request.JsonData)
	//fmt.Println("request.AccountID: " + request.AccountId)
	//fmt.Println("request.JsonData: " + request.JsonData)

	sessionState, err := services.CreateSessionState(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, sessionState, "")
}

func (handler *stateSessionsHandler) delete(c *gin.Context) {
	accountId := c.Param("accountId")
	stateType := c.Param("stateType")

	//DeleteSessionStateForAccountIdStateType

	iAccountId, _ := strconv.ParseUint(accountId, 10, 64)
	iStateType, _ := strconv.ParseUint(stateType, 10, 64)

	services.DeleteSessionState(iAccountId, iStateType)

	dtos.JsonResponse(c, http.StatusOK, nil, "")

}
