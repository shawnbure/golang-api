package handlers

import (
	"fmt"
	"net/http"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/proxy/middleware"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseSessionStatesEndpoint   = "/session-states"
	sessionStatesCreateEndpoint = "/create"
	sessionStatesDeleteEndpoint = "/delete"
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

	err := c.BindJSON(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	fmt.Println("request.JsonData: " + request.JsonData)
	//fmt.Println("request.AccountID: " + request.AccountId)
	//fmt.Println("request.JsonData: " + request.JsonData)

	sessionState, err := services.CreateSessionState(&request)
	if err != nil {
		fmt.Println("error: " + err.Error())
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, sessionState, "")
}

func (handler *stateSessionsHandler) delete(c *gin.Context) {
	var request services.DeleteSessionStateRequest

	err := c.BindJSON(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	fmt.Println("proxy handler > sessionState: delete")
	fmt.Println(request.AccountId)
	fmt.Println(request.StateType)

	strResult, err := services.DeleteSessionState(&request)
	if err != nil {
		fmt.Println("error: " + err.Error())
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, strResult, "")

}
