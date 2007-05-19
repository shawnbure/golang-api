package handlers

import (
	"net/http"
	"time"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/services"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseAccountsEndpoint         = "/accounts"
	accountByUserAddressEndpoint = "/:userAddress"
	accountProfileEndpoint       = "/images/profile/:userAddress"
	accountCoverEndpoint         = "/images/cover/:userAddress"
)

type accountsHandler struct {
}

func NewAccountsHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &accountsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: accountByUserAddressEndpoint, HandlerFunc: handler.get},
		{Method: http.MethodPost, Path: accountByUserAddressEndpoint, HandlerFunc: handler.set},

		{Method: http.MethodGet, Path: accountProfileEndpoint, HandlerFunc: handler.getAccountProfile},
		{Method: http.MethodPost, Path: accountProfileEndpoint, HandlerFunc: handler.setAccountProfile},

		{Method: http.MethodGet, Path: accountCoverEndpoint, HandlerFunc: handler.getAccountCover},
		{Method: http.MethodPost, Path: accountCoverEndpoint, HandlerFunc: handler.setAccountCover},
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

	data.JsonResponse(c, http.StatusOK, account, "")
}

func (handler *accountsHandler) getAccountProfile(c *gin.Context) {
	userAddress := c.Param("userAddress")

	image, err := services.GetAccountProfileImage(userAddress)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, &image, "")
}

func (handler *accountsHandler) setAccountProfile(c *gin.Context) {
	var imageBase64 string
	userAddress := c.Param("userAddress")

	err := c.Bind(&imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if jwtAddress != userAddress {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	err = services.SetAccountProfileImage(userAddress, &imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, "", "")
}

func (handler *accountsHandler) getAccountCover(c *gin.Context) {
	userAddress := c.Param("userAddress")

	image, err := services.GetAccountCoverImage(userAddress)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, &image, "")
}

func (handler *accountsHandler) setAccountCover(c *gin.Context) {
	var imageBase64 string
	userAddress := c.Param("userAddress")

	err := c.Bind(&imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if jwtAddress != userAddress {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	err = services.SetAccountCoverImage(userAddress, &imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, "", "")
}
