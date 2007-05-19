package handlers

import (
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/services"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseAccountsEndpoint   = "/accounts"
	createAccountEndpoint  = "/create"
	accountByIdEndpoint    = "/:walletAddress"
	accountTokensEndpoint  = "/:walletAddress/tokens/:offset/:limit"
	accountProfileEndpoint = "/:walletAddress/profile"
	accountCoverEndpoint   = "/:walletAddress/cover"
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

		{Method: http.MethodPost, Path: createAccountEndpoint, HandlerFunc: handler.create},
		{Method: http.MethodGet, Path: accountTokensEndpoint, HandlerFunc: handler.getAccountTokens},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseAccountsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Get account by account walletAddress
// @Description Retrieves an account by walletAddress
// @Tags accounts
// @Accept json
// @Produce json
// @Param walletAddress path string true "wallet address"
// @Success 200 {object} entities.Account
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /accounts/{accountId} [get]
func (handler *accountsHandler) get(c *gin.Context) {
	walletAddress := c.Param("walletAddress")

	cacheInfo, err := services.GetOrAddAccountCacheInfo(walletAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	account, err := storage.GetAccountById(cacheInfo.AccountId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, account, "")
}

// @Summary Set account information
// @Description Sets an account settable information
// @Tags accounts
// @Accept json
// @Produce json
// @Param walletAddress path string true "wallet address"
// @Param setAccountRequest body services.SetAccountRequest true "account info"
// @Success 200 {object} entities.Account
// @Failure 400 {object} dtos.ApiResponse
// @Failure 401 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /accounts/{walletAddress} [post]
func (handler *accountsHandler) set(c *gin.Context) {
	var request services.SetAccountRequest
	walletAddress := c.Param("walletAddress")

	err := c.Bind(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, "cannot bind request")
		return
	}

	cacheInfo, err := services.GetOrAddAccountCacheInfo(walletAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	account, err := storage.GetAccountById(cacheInfo.AccountId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if account.Address != jwtAddress {
		dtos.JsonResponse(c, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	err = services.UpdateAccount(account, &request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, "could not get price")
		return
	}

	dtos.JsonResponse(c, http.StatusOK, account, "")
}

// @Summary Creates an account
// @Description Creates an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param createAccountRequest body services.CreateAccountRequest true "account info"
// @Success 200 {object} entities.Account
// @Failure 400 {object} dtos.ApiResponse
// @Failure 401 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /accounts/{address} [post]
func (handler *accountsHandler) create(c *gin.Context) {
	var request services.CreateAccountRequest

	err := c.Bind(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, "cannot bind request")
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if request.Address != jwtAddress {
		dtos.JsonResponse(c, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	_, err = storage.GetAccountByAddress(request.Address)
	if err == nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, "account already exists for address")
		return
	}

	account, err := services.CreateAccount(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, "could not get price")
		return
	}

	dtos.JsonResponse(c, http.StatusOK, account, "")
}

// @Summary Get account profile image
// @Description Retrieves an account profile image. It will be sent as base64 encoding (sdt, raw) of its byte representation.
// @Tags accounts
// @Accept json
// @Produce json
// @Param walletAddress path uint64 true "wallet address"
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /accounts/{walletAddress}/profile [get]
func (handler *accountsHandler) getAccountProfile(c *gin.Context) {
	walletAddress := c.Param("walletAddress")

	cacheInfo, err := services.GetOrAddAccountCacheInfo(walletAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	image, err := storage.GetAccountProfileImageByAccountId(cacheInfo.AccountId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, &image, "")
}

// @Summary Set account profile image
// @Description Expects base64 encoding (sdt, raw) of the image representation. Returns empty string. Max size of byte array is 512KB.
// @Tags accounts
// @Accept json
// @Produce json
// @Param walletAddress path uint64 true "wallet address"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 401 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /accounts/{walletAddress}/profile [post]
func (handler *accountsHandler) setAccountProfile(c *gin.Context) {
	var imageBase64 string
	walletAddress := c.Param("walletAddress")

	err := c.Bind(&imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if jwtAddress != walletAddress {
		dtos.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	cacheInfo, err := services.GetOrAddAccountCacheInfo(walletAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	err = services.SetAccountProfileImage(cacheInfo.AccountId, &imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, "", "")
}

// @Summary Get account cover image
// @Description Retrieves an account cover image. It will be sent as base64 encoding (sdt, raw) of its byte representation.
// @Tags accounts
// @Accept json
// @Produce json
// @Param walletAddress path uint64 true "wallet address"
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /accounts/{walletAddress}/cover [get]
func (handler *accountsHandler) getAccountCover(c *gin.Context) {
	walletAddress := c.Param("walletAddress")

	cacheInfo, err := services.GetOrAddAccountCacheInfo(walletAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	image, err := storage.GetAccountCoverImageByAccountId(cacheInfo.AccountId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, &image, "")
}

// @Summary Set account cover image
// @Description Expects base64 encoding (sdt, raw) of the image representation. Returns empty string. Max size of byte array is 1MB.
// @Tags accounts
// @Accept json
// @Produce json
// @Param walletAddress path uint64 true "wallet address"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 401 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /accounts/{walletAddress}/cover [post]
func (handler *accountsHandler) setAccountCover(c *gin.Context) {
	var imageBase64 string
	walletAddress := c.Param("walletAddress")

	err := c.Bind(&imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if jwtAddress != walletAddress {
		dtos.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	cacheInfo, err := services.GetOrAddAccountCacheInfo(walletAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	err = services.SetAccountCoverImage(cacheInfo.AccountId, &imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, "", "")
}

// @Summary Gets tokens for an account.
// @Description Retrieves a list of tokens. Unsorted.
// @Tags accounts
// @Accept json
// @Produce json
// @Param walletAddress path uint64 true "wallet address"
// @Param offset path int true "offset"
// @Param limit path int true "limit"
// @Success 200 {object} []entities.Token
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /accounts/{walletAddress}/tokens/{offset}/{limit} [get]
func (handler *accountsHandler) getAccountTokens(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")
	walletAddress := c.Param("walletAddress")

	cacheInfo, err := services.GetOrAddAccountCacheInfo(walletAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	tokens, err := storage.GetTokensByOwnerIdWithOffsetLimit(cacheInfo.AccountId, offset, limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, tokens, "")
}
