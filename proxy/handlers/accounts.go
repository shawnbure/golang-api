package handlers

import (
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/services"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	baseAccountsEndpoint = "/accounts"
	getAccountEndpoint   = "/get/:userAddress"
	setAccountEndpoint   = "/set/:userAddress"
)

type accountsHandler struct {
}

func NewAccountsHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &accountsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: getAccountEndpoint, HandlerFunc: handler.get},
		{Method: http.MethodPost, Path: setAccountEndpoint, HandlerFunc: handler.set},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseAccountsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *accountsHandler) get(c *gin.Context) {
	userAddress := c.Param("userAddress")

	account, err := storage.GetAccountByAddress(userAddress)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, "could not get price")
		return
	}

	data.JsonResponse(c, http.StatusOK, account, "")
}

func (handler *accountsHandler) set(c *gin.Context) {
	var request services.SetAccountRequest
	userAddress := c.Param("userAddress")

	err := c.Bind(&request)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, "cannot bind request")
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if userAddress != jwtAddress {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	account := data.Account{
		Address:       userAddress,
		Name:          request.Name,
		Description:   request.Description,
		Website:       request.Website,
		TwitterLink:   request.TwitterLink,
		InstagramLink: request.InstagramLink,
		CreatedAt:     uint64(time.Now().Unix()),
	}
	err = services.AddOrUpdateAccount(&account)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, "could not get price")
		return
	}

	data.JsonResponse(c, http.StatusOK, "", "")
}
