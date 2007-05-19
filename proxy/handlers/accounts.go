package handlers

import (
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/services"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseAccountsEndpoint   = "/accounts"
	accountByIdEndpoint    = "/:accountId"
	accountAssetsEndpoint  = "/:accountId/assets/:offset/:limit"
	accountProfileEndpoint = "/:accountId/profile"
	accountCoverEndpoint   = "/:accountId/cover"

	accountByAddressEndpoint = "/find/:accountAddress"
	createAccountEndpoint    = "/create"
)

type accountsHandler struct {
}

func NewAccountsHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &accountsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: accountByIdEndpoint, HandlerFunc: handler.get},
		{Method: http.MethodPost, Path: accountByIdEndpoint, HandlerFunc: handler.set},

		{Method: http.MethodGet, Path: accountProfileEndpoint, HandlerFunc: handler.getAccountProfile},
		{Method: http.MethodPost, Path: accountProfileEndpoint, HandlerFunc: handler.setAccountProfile},

		{Method: http.MethodGet, Path: accountCoverEndpoint, HandlerFunc: handler.getAccountCover},
		{Method: http.MethodPost, Path: accountCoverEndpoint, HandlerFunc: handler.setAccountCover},

		{Method: http.MethodGet, Path: accountByAddressEndpoint, HandlerFunc: handler.getByAddress},
		{Method: http.MethodPost, Path: createAccountEndpoint, HandlerFunc: handler.create},
		{Method: http.MethodGet, Path: accountAssetsEndpoint, HandlerFunc: handler.getAccountAssets},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseAccountsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Get account by account id
// @Description Retrieves an account by id
// @Tags accounts
// @Accept json
// @Produce json
// @Param accountId path string true "account id"
// @Success 200 {object} data.Account
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /accounts/{accountId} [get]
func (handler *accountsHandler) get(c *gin.Context) {
	accountIdString := c.Param("accountId")

	accountId, err := strconv.ParseUint(accountIdString, 10, 16)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	account, err := storage.GetAccountById(accountId)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, "could not get price")
		return
	}

	data.JsonResponse(c, http.StatusOK, account, "")
}

// @Summary Get account by address
// @Description Retrieves an account by address. Useful for login.
// @Tags accounts
// @Accept json
// @Produce json
// @Param accountAddress path string true "account address"
// @Success 200 {object} data.Account
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /accounts/find/{accountAddress} [get]
func (handler *accountsHandler) getByAddress(c *gin.Context) {
	accountAddress := c.Param("accountAddress")

	account, err := storage.GetAccountByAddress(accountAddress)
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
// @Param accountId path string true "account id"
// @Param setAccountRequest body services.SetAccountRequest true "account info"
// @Success 200 {object} data.Account
// @Failure 400 {object} data.ApiResponse
// @Failure 401 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /accounts/{accountId} [post]
func (handler *accountsHandler) set(c *gin.Context) {
	var request services.SetAccountRequest
	accountIdString := c.Param("accountId")

	err := c.Bind(&request)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, "cannot bind request")
		return
	}

	accountId, err := strconv.ParseUint(accountIdString, 10, 16)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	account, err := storage.GetAccountById(accountId)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if account.Address != jwtAddress {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	err = services.UpdateAccount(account, &request)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, "could not get price")
		return
	}

	data.JsonResponse(c, http.StatusOK, account, "")
}

// @Summary Creates an account
// @Description Creates an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param createAccountRequest body services.CreateAccountRequest true "account info"
// @Success 200 {object} data.Account
// @Failure 400 {object} data.ApiResponse
// @Failure 401 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /accounts/{address} [post]
func (handler *accountsHandler) create(c *gin.Context) {
	var request services.CreateAccountRequest

	err := c.Bind(&request)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, "cannot bind request")
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if request.Address != jwtAddress {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	_, err = storage.GetAccountByAddress(request.Address)
	if err == nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, "account already exists for address")
		return
	}

	account, err := services.CreateAccount(&request)
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
// @Param accountId path uint64 true "account id"
// @Success 200 {object} string
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /accounts/{accountId}/profile [get]
func (handler *accountsHandler) getAccountProfile(c *gin.Context) {
	accountIdString := c.Param("accountId")

	accountId, err := strconv.ParseUint(accountIdString, 10, 16)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	image, err := storage.GetAccountProfileImageByAccountId(accountId)
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
// @Param accountId path uint64 true "account id"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} data.ApiResponse
// @Failure 401 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /accounts/{accountId}/profile [post]
func (handler *accountsHandler) setAccountProfile(c *gin.Context) {
	var imageBase64 string
	accountIdString := c.Param("accountId")

	err := c.Bind(&imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	accountId, err := strconv.ParseUint(accountIdString, 10, 16)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	account, err := storage.GetAccountById(accountId)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	if jwtAddress != account.Address {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	err = services.SetAccountProfileImage(accountId, &imageBase64)
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
// @Param accountId path uint64 true "account id"
// @Success 200 {object} string
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /accounts/{accountId}/cover [get]
func (handler *accountsHandler) getAccountCover(c *gin.Context) {
	accountIdString := c.Param("accountId")

	accountId, err := strconv.ParseUint(accountIdString, 10, 16)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	image, err := storage.GetAccountCoverImageByAccountId(accountId)
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
// @Param accountId path uint64 true "account id"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} data.ApiResponse
// @Failure 401 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /accounts/{accountId}/cover [post]
func (handler *accountsHandler) setAccountCover(c *gin.Context) {
	var imageBase64 string
	accountIdString := c.Param("accountId")

	err := c.Bind(&imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	accountId, err := strconv.ParseUint(accountIdString, 10, 16)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	account, err := storage.GetAccountById(accountId)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	if jwtAddress != account.Address {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	err = services.SetAccountCoverImage(accountId, &imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, "", "")
}

// @Summary Gets assets for an account.
// @Description Retrieves a list of assets. Unsorted.
// @Tags accounts
// @Accept json
// @Produce json
// @Param accountId path uint64 true "account id"
// @Param offset path int true "offset"
// @Param limit path int true "limit"
// @Success 200 {object} []data.Asset
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /accounts/{accountId}/assets/{offset}/{limit} [get]
func (handler *accountsHandler) getAccountAssets(c *gin.Context) {
	accountIdString := c.Param("accountId")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	accountId, err := strconv.ParseUint(accountIdString, 10, 16)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	assets, err := storage.GetAssetsByOwnerIdWithOffsetLimit(accountId, offset, limit)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, assets, "")
}
