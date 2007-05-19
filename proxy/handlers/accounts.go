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
	accountProfileEndpoint       = "/:userAddress/profile"
	accountCoverEndpoint         = "/:userAddress/cover"
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

// @Summary Get account by user address
// @Description Retrieves an account by an elrond user address (erd1...)
// @Tags accounts
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Success 200 {object} data.Account
// @Failure 404 {object} data.ApiResponse
// @Router /accounts/{userAddress} [get]
func (handler *accountsHandler) get(c *gin.Context) {
	userAddress := c.Param("userAddress")

	account, err := storage.GetAccountByAddress(userAddress)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, "could not get price")
		return
	}

	data.JsonResponse(c, http.StatusOK, account, "")
}

// @Summary Set account information
// @Description Sets an account settable information
// @Tags accounts
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Param setAccountRequest body services.SetAccountRequest true "account info"
// @Success 200 {object} data.Account
// @Failure 400 {object} data.ApiResponse
// @Failure 401 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /accounts/{userAddress} [post]
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

// @Summary Get account profile image
// @Description Retrieves an account profile image. It will be sent as base64 encoding (sdt, raw) of its byte representation.
// @Tags accounts
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Success 200 {object} string
// @Failure 404 {object} data.ApiResponse
// @Router /accounts/{userAddress}/profile [get]
func (handler *accountsHandler) getAccountProfile(c *gin.Context) {
	userAddress := c.Param("userAddress")

	image, err := services.GetAccountProfileImage(userAddress)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, &image, "")
}

// @Summary Set account profile image
// @Description Expects base64 encoding (sdt, raw) of the image representation. Returns empty string. Max size of byte array is 512KB.
// @Tags accounts
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} data.ApiResponse
// @Failure 401 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /accounts/{userAddress}/profile [post]
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

// @Summary Get account cover image
// @Description Retrieves an account cover image. It will be sent as base64 encoding (sdt, raw) of its byte representation.
// @Tags accounts
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Success 200 {object} string
// @Failure 404 {object} data.ApiResponse
// @Router /accounts/{userAddress}/cover [get]
func (handler *accountsHandler) getAccountCover(c *gin.Context) {
	userAddress := c.Param("userAddress")

	image, err := services.GetAccountCoverImage(userAddress)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, &image, "")
}

// @Summary Set account cover image
// @Description Expects base64 encoding (sdt, raw) of the image representation. Returns empty string. Max size of byte array is 1MB.
// @Tags accounts
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} data.ApiResponse
// @Failure 401 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /accounts/{userAddress}/cover [post]
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
