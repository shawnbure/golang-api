package handlers

import (
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/services"
	"github.com/erdsea/erdsea-api/stats/collstats"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseCollectionsEndpoint      = "/collections"
	collectionByNameEndpoint     = "/:collectionId"
	collectionListEndpoint       = "/list/:offset/:limit"
	collectionCreateEndpoint     = "/create"
	collectionTokensEndpoint     = "/:collectionId/tokens/:offset/:limit"
	collectionProfileEndpoint    = "/:collectionId/profile/"
	collectionCoverEndpoint      = "/:collectionId/cover"
)

type collectionsHandler struct {
	blockchainCfg config.BlockchainConfig
}

func NewCollectionsHandler(groupHandler *groupHandler, authCfg config.AuthConfig, blockchainCfg config.BlockchainConfig) {
	handler := &collectionsHandler{blockchainCfg: blockchainCfg}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: collectionListEndpoint, HandlerFunc: handler.getList},
		{Method: http.MethodGet, Path: collectionByNameEndpoint, HandlerFunc: handler.get},
		{Method: http.MethodPost, Path: collectionByNameEndpoint, HandlerFunc: handler.set},

		{Method: http.MethodPost, Path: collectionCreateEndpoint, HandlerFunc: handler.create},
		{Method: http.MethodGet, Path: collectionTokensEndpoint, HandlerFunc: handler.getTokens},

		{Method: http.MethodGet, Path: collectionProfileEndpoint, HandlerFunc: handler.getCollectionProfile},
		{Method: http.MethodPost, Path: collectionProfileEndpoint, HandlerFunc: handler.setCollectionProfile},

		{Method: http.MethodGet, Path: collectionCoverEndpoint, HandlerFunc: handler.getCollectionCover},
		{Method: http.MethodPost, Path: collectionCoverEndpoint, HandlerFunc: handler.setCollectionCover},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseCollectionsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Gets collections.
// @Description Retrieves a list of collections. Unsorted.
// @Tags collections
// @Accept json
// @Produce json
// @Param offset path int true "offset"
// @Param limit path int true "limit"
// @Success 200 {object} []entities.Collection
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /collections/list/{offset}/{limit} [get]
func (handler *collectionsHandler) getList(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

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

	collections, err := storage.GetCollectionsWithOffsetLimit(offset, limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, collections, "")
}

// @Summary Gets collection by collection id.
// @Description Retrieves a collection by id.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionId path string true "collection id"
// @Success 200 {object} entities.Collection
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /collections/{collectionId} [get]
func (handler *collectionsHandler) get(c *gin.Context) {
	tokenId := c.Param("collectionId")

	cacheInfo, err := collstats.GetCollection(tokenId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	collection, err := storage.GetCollectionById(cacheInfo.CollectionId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, collection, "")
}

// @Summary Set collection info.
// @Description Sets info for a collection.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionId path string true "collection id"
// @Param updateCollectionRequest body services.UpdateCollectionRequest true "collection info"
// @Success 200 {object} entities.Collection
// @Failure 401 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /collections/{collectionId} [post]
func (handler *collectionsHandler) set(c *gin.Context) {
	var request services.UpdateCollectionRequest
	tokenId := c.Param("collectionId")

	err := c.Bind(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	cacheInfo, err := collstats.GetCollection(tokenId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	collection, err := storage.GetCollectionById(cacheInfo.CollectionId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	creator, err := storage.GetAccountById(collection.CreatorID)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if creator.Address != jwtAddress {
		dtos.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	err = services.UpdateCollection(collection, &request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, collection, "")
}

// @Summary Creates a collection.
// @Description Creates a collection with given info.
// @Tags collections
// @Accept json
// @Produce json
// @Param createCollectionRequest body services.CreateCollectionRequest true "collection info"
// @Success 200 {object} entities.Collection
// @Failure 400 {object} dtos.ApiResponse
// @Failure 401 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /collections/create [post]
func (handler *collectionsHandler) create(c *gin.Context) {
	var request services.CreateCollectionRequest

	err := c.Bind(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if jwtAddress != request.UserAddress {
		dtos.JsonResponse(c, http.StatusUnauthorized, nil, "jwt and request addresses differ")
		return
	}

	collection, err := services.CreateCollection(&request, handler.blockchainCfg.ProxyUrl)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, collection, "")
}

// @Summary Get collection tokens.
// @Description Retrieves the tokens of a collection. Unsorted.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionId path string true "collection id"
// @Param offset path int true "offset"
// @Param limit path int true "limit"
// @Success 200 {object} []entities.Token
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /collections/{collectionId}/tokens/{offset}/{limit} [get]
func (handler *collectionsHandler) getTokens(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")
	filters := c.QueryMap("filters")
	tokenId := c.Param("collectionId")

	cacheInfo, err := collstats.GetCollection(tokenId)
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

	tokens, err := storage.GetTokensByCollectionIdWithOffsetLimit(cacheInfo.CollectionId, offset, limit, filters)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, tokens, "")
}

// @Summary Get collection profile image
// @Description Retrieves a collection cover image. It will be sent as base64 encoding (sdt, raw) of its byte representation.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionId path string true "collection id"
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /collections/{collectionId}/profile [get]
func (handler *collectionsHandler) getCollectionProfile(c *gin.Context) {
	tokenId := c.Param("collectionId")

	cacheInfo, err := collstats.GetCollection(tokenId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	image, err := storage.GetCollectionProfileImageByCollectionId(cacheInfo.CollectionId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, &image, "")
}

// @Summary Set collection profile image
// @Description Expects base64 encoding (sdt, raw) of the image representation. Returns empty string. Max size of byte array is 1MB.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionId path string true "collection id"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /collections/{collectionId}/profile [post]
func (handler *collectionsHandler) setCollectionProfile(c *gin.Context) {
	var imageBase64 string
	tokenId := c.Param("collectionId")

	err := c.Bind(&imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	cacheInfo, err := collstats.GetCollection(tokenId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	collection, err := storage.GetCollectionById(cacheInfo.CollectionId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	account, err := storage.GetAccountById(collection.CreatorID)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if account.Address != jwtAddress {
		dtos.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	err = services.SetCollectionProfileImage(cacheInfo.CollectionId, &imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, "", "")
}

// @Summary Get collection cover image
// @Description Retrieves a collection cover image. It will be sent as base64 encoding (sdt, raw) of its byte representation.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionId path string true "collection id"
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /collections/{collectionId}/cover [get]
func (handler *collectionsHandler) getCollectionCover(c *gin.Context) {
	tokenId := c.Param("collectionId")

	cacheInfo, err := collstats.GetCollection(tokenId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	image, err := storage.GetCollectionCoverImageByCollectionId(cacheInfo.CollectionId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, &image, "")
}

// @Summary Set collection cover image
// @Description Expects base64 encoding (sdt, raw) of the image representation. Returns empty string. Max size of byte array is 1MB.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionId path string true "collection id"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} dtos.ApiResponse
// @Failure 401 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /collections/{collectionId}/cover [post]
func (handler *collectionsHandler) setCollectionCover(c *gin.Context) {
	var imageBase64 string
	tokenId := c.Param("collectionId")

	err := c.Bind(&imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	cacheInfo, err := collstats.GetCollection(tokenId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	collection, err := storage.GetCollectionById(cacheInfo.CollectionId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	account, err := storage.GetAccountById(collection.CreatorID)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if account.Address != jwtAddress {
		dtos.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	err = services.SetCollectionCoverImage(cacheInfo.CollectionId, &imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, "", "")
}
