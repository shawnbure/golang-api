package handlers

import (
	"bytes"
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
	baseAccountsEndpoint       = "/accounts"
	accountByIdEndpoint        = "/:walletAddress"
	accountTokensEndpoint      = "/:walletAddress/tokens/:offset/:limit"
	accountCollectionsEndpoint = "/:walletAddress/collections/:offset/:limit"
	accountProfileEndpoint     = "/:walletAddress/profile"
	accountCoverEndpoint       = "/:walletAddress/cover"
	imageEndpoint              = "/image/:filename"
)

type accountsHandler struct {
}

func NewAccountsHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &accountsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodPost, Path: accountByIdEndpoint, HandlerFunc: handler.set},
		{Method: http.MethodPost, Path: accountProfileEndpoint, HandlerFunc: handler.setAccountProfile},
		{Method: http.MethodPost, Path: accountCoverEndpoint, HandlerFunc: handler.setAccountCover},
	}
	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseAccountsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}
	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)

	publicEndpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: accountByIdEndpoint, HandlerFunc: handler.get},
		{Method: http.MethodGet, Path: accountTokensEndpoint, HandlerFunc: handler.getAccountTokens},
		{Method: http.MethodGet, Path: accountCollectionsEndpoint, HandlerFunc: handler.getAccountCollections},
	}
	publicEndpointGroupHandler := EndpointGroupHandler{
		Root:             baseAccountsEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: publicEndpoints,
	}
	groupHandler.AddEndpointGroupHandler(publicEndpointGroupHandler)
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
// @Router /accounts/{walletAddress} [get]
func (h *accountsHandler) get(c *gin.Context) {
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
func (h *accountsHandler) set(c *gin.Context) {
	var request services.SetAccountRequest
	walletAddress := c.Param("walletAddress")

	err := c.BindJSON(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, "cannot bind request")
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if walletAddress != jwtAddress {
		dtos.JsonResponse(c, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	var innerErr error
	account, err := storage.GetAccountByAddress(walletAddress)
	if err != nil {
		account, innerErr = services.CreateAccount(walletAddress, &request)
	} else {
		innerErr = services.UpdateAccount(account, &request)
	}

	if innerErr != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, innerErr.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, account, "")
}

// @Summary Set account profile image
// @Description Expects base64 std encoding of the image representation. Returns empty string. Max size of byte array is 512KB.
// @Tags accounts
// @Accept json
// @Produce json
// @Param walletAddress path string true "wallet address"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 401 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /accounts/{walletAddress}/profile [post]
func (h *accountsHandler) setAccountProfile(c *gin.Context) {
	walletAddress := c.Param("walletAddress")

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	imageBase64 := buf.String()
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

	link, err := services.SetAccountProfileImage(walletAddress, cacheInfo.AccountId, &imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, link, "")
}

// @Summary Set account cover image
// @Description Expects base64 std encoding of the image representation. Returns empty string. Max size of byte array is 1MB.
// @Tags accounts
// @Accept json
// @Produce json
// @Param walletAddress path string true "wallet address"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 401 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /accounts/{walletAddress}/cover [post]
func (h *accountsHandler) setAccountCover(c *gin.Context) {
	walletAddress := c.Param("walletAddress")

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	imageBase64 := buf.String()
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

	link, err := services.SetAccountCoverImage(walletAddress, cacheInfo.AccountId, &imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, link, "")
}

// @Summary Gets tokens for an account.
// @Description Retrieves a list of tokens. Unsorted.
// @Tags accounts
// @Accept json
// @Produce json
// @Param walletAddress path string true "wallet address"
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []dtos.OwnedTokenDto
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /accounts/{walletAddress}/tokens/{offset}/{limit} [get]
func (h *accountsHandler) getAccountTokens(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")
	walletAddress := c.Param("walletAddress")

	cacheInfo, err := services.GetOrAddAccountCacheInfo(walletAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	offset, err := strconv.ParseUint(offsetStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	err = ValidateLimit(limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	tokens, err := storage.GetTokensByOwnerIdWithOffsetLimit(cacheInfo.AccountId, int(offset), int(limit))
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	ownedTokens := services.ConstructOwnedTokensFromTokens(tokens)
	dtos.JsonResponse(c, http.StatusOK, ownedTokens, "")
}

// @Summary Gets collections for an account.
// @Description Retrieves a list of collections. Unsorted.
// @Tags accounts
// @Accept json
// @Produce json
// @Param walletAddress path string true "wallet address"
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.Collection
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /accounts/{walletAddress}/collections/{offset}/{limit} [get]
func (h *accountsHandler) getAccountCollections(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")
	walletAddress := c.Param("walletAddress")

	cacheInfo, err := services.GetOrAddAccountCacheInfo(walletAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	offset, err := strconv.ParseUint(offsetStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	err = ValidateLimit(limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	collections, err := storage.GetCollectionsByCreatorIdWithOffsetLimit(cacheInfo.AccountId, int(offset), int(limit))
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, collections, "")
}
