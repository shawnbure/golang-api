package handlers

import (
	"encoding/json"
	"github.com/erdsea/erdsea-api/stats/collstats"
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
	baseTokensEndpoint               = "/tokens"
	tokenByTokenIdAndNonceEndpoint   = "/:tokenId/:nonce"
	availableTokensEndpoint          = "/available"
	offersForTokenIdAndNonceEndpoint = "/:tokenId/:nonce/offers/:offset/:limit"
	bidsForTokenIdAndNonceEndpoint   = "/:tokenId/:nonce/bids/:offset/:limit"
	refreshTokenMetadataEndpoint     = "/:tokenId/:nonce/refresh"
	tokenMetadataRelayEndpoint       = "/metadata/relay"
	tokensListMetadataEndpoint       = "/list/:offset/:limit"
)

type TokenListQueryBody struct {
	SortRules map[string]string `json:"sortRules"`
}

type tokensHandler struct {
	blockchainConfig config.BlockchainConfig
}

func NewTokensHandler(groupHandler *groupHandler, cfg config.BlockchainConfig) {
	handler := &tokensHandler{
		blockchainConfig: cfg,
	}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: tokenByTokenIdAndNonceEndpoint, HandlerFunc: handler.getByTokenIdAndNonce},
		{Method: http.MethodPost, Path: availableTokensEndpoint, HandlerFunc: handler.getAvailableTokens},
		{Method: http.MethodGet, Path: offersForTokenIdAndNonceEndpoint, HandlerFunc: handler.getOffers},
		{Method: http.MethodGet, Path: bidsForTokenIdAndNonceEndpoint, HandlerFunc: handler.getBids},
		{Method: http.MethodGet, Path: tokenMetadataRelayEndpoint, HandlerFunc: handler.relayMetadataResponse},
		{Method: http.MethodPost, Path: refreshTokenMetadataEndpoint, HandlerFunc: handler.refresh},
		{Method: http.MethodPost, Path: tokensListMetadataEndpoint, HandlerFunc: handler.getList},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseTokensEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Get token by id and nonce
// @Description Retrieves a token by tokenId and nonce
// @Tags tokens
// @Accept json
// @Produce json
// @Param tokenId path string true "token id"
// @Param nonce path int true "token nonce"
// @Success 200 {object} dtos.ExtendedTokenDto
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /tokens/{tokenId}/{nonce} [get]
func (handler *tokensHandler) getByTokenIdAndNonce(c *gin.Context) {
	tokenId := c.Param("tokenId")
	nonceString := c.Param("nonce")

	nonce, err := strconv.ParseUint(nonceString, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	tokenDto, err := services.GetExtendedTokenData(tokenId, nonce)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, tokenDto, "")
}

// @Summary Get available tokens
// @Description Get available tokens and some collection info
// @Tags tokens
// @Accept json
// @Produce json
// @Param availableTokensRequest body services.AvailableTokensRequest true "request"
// @Success 200 {object} services.AvailableTokensResponse
// @Failure 400 {object} dtos.ApiResponse
// @Router /tokens/available [get]
func (handler *tokensHandler) getAvailableTokens(c *gin.Context) {
	var request services.AvailableTokensRequest

	err := c.Bind(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	response, err := services.GetAvailableTokens(request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, response, "")
}

// @Summary Get offers for token
// @Description Retrieves offers for a token (identified by tokenId and nonce)
// @Tags tokens
// @Accept json
// @Produce json
// @Param tokenId path string true "token id"
// @Param nonce path int true "token nonce"
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []dtos.OfferDto
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /tokens/{tokenId}/{nonce}/offers/{offset}/{limit} [get]
func (handler *tokensHandler) getOffers(c *gin.Context) {
	tokenId := c.Param("tokenId")
	nonceString := c.Param("nonce")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	nonce, err := strconv.ParseUint(nonceString, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
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

	tokenCacheInfo, err := services.GetOrAddTokenCacheInfo(tokenId, nonce)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	offers, err := storage.GetOffersForTokenWithOffsetLimit(tokenCacheInfo.TokenDbId, int(offset), int(limit))
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	offersDtos := services.MakeOfferDtos(offers)
	dtos.JsonResponse(c, http.StatusOK, offersDtos, "")
}

// @Summary Get bids for token
// @Description Retrieves bids for a token (identified by tokenId and nonce)
// @Tags tokens
// @Accept json
// @Produce json
// @Param tokenId path string true "token id"
// @Param nonce path int true "token nonce"
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []dtos.BidDto
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /tokens/{tokenId}/{nonce}/bids/{offset}/{limit} [get]
func (handler *tokensHandler) getBids(c *gin.Context) {
	tokenId := c.Param("tokenId")
	nonceString := c.Param("nonce")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	nonce, err := strconv.ParseUint(nonceString, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
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

	tokenCacheInfo, err := services.GetOrAddTokenCacheInfo(tokenId, nonce)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	bids, err := storage.GetBidsForTokenWithOffsetLimit(tokenCacheInfo.TokenDbId, int(offset), int(limit))
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	bidsDtos := services.MakeBidDtos(bids)
	dtos.JsonResponse(c, http.StatusOK, bidsDtos, "")
}

// @Summary Gets metadata link response. Cached.
// @Description Make request with ?url=link
// @Tags tokens
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /tokens/metadata/relay [get]
func (handler *tokensHandler) relayMetadataResponse(c *gin.Context) {
	url := c.Query("url")

	responseBytes, err := services.TryGetResponseCached(url)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	var metadata dtos.MetadataLinkResponse
	err = json.Unmarshal([]byte(responseBytes), &metadata)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, metadata, "")
}

// @Summary Tries to refresh token metadata link and attributes.
// @Description Returns attributes directly stored inside token (not OS format). Check then before and after. If modified, reload the page maybe?
// @Tags tokens
// @Accept json
// @Produce json
// @Param tokenId path string true "token id"
// @Param nonce path int true "token nonce"
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /tokens/{tokenId}/{nonce}/refresh [post]
func (handler *tokensHandler) refresh(c *gin.Context) {
	tokenId := c.Param("tokenId")
	nonceString := c.Param("nonce")

	nonce, err := strconv.ParseUint(nonceString, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	collectionCacheInfo, err := collstats.GetOrAddCollectionCacheInfo(tokenId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	proxy := handler.blockchainConfig.ProxyUrl
	jwtAddress := c.GetString(middleware.AddressKey)
	collectionId := collectionCacheInfo.CollectionId
	marketplace := handler.blockchainConfig.MarketplaceAddress
	metadata, err := services.AddOrRefreshToken(tokenId, nonce, collectionId, jwtAddress, proxy, marketplace)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, "")
		return
	}

	dtos.JsonResponse(c, http.StatusOK, metadata, "")
}

// @Summary Gets tokens.
// @Description Retrieves a list of tokens.
// @Tags tokens
// @Accept json
// @Produce json
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Param query body TokenListQueryBody true "sort rules"
// @Success 200 {object} []entities.Token
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /tokens/list/{offset}/{limit} [post]
func (handler *tokensHandler) getList(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	var queries TokenListQueryBody
	err := c.BindJSON(&queries)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	sortRules := queries.SortRules

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

	if sortRules == nil {
		sortRules = make(map[string]string, 2)
	}
	if len(sortRules) == 0 {
		sortRules["criteria"] = "created_at"
		sortRules["mode"] = "desc"
	}

	acceptedCriteria := map[string]bool{"price_nominal": true, "created_at": true}
	err = testInputSortParams(sortRules, acceptedCriteria)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	tokens, err := storage.GetTokensWithOffsetLimit(int(offset), int(limit), sortRules)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, tokens, "")
}
