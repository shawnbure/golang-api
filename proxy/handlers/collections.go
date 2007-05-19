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
	baseCollectionsEndpoint      = "/collections"
	collectionByNameEndpoint     = "/:collectionName"
	collectionListEndpoint       = "/list/:offset/:limit"
	collectionCreateEndpoint     = "/create"
	collectionStatisticsEndpoint = "/:collectionName/statistics"
	collectionAssetsEndpoint     = "/:collectionName/assets/:offset/:limit"
	collectionProfileEndpoint    = "/:collectionName/profile/"
	collectionCoverEndpoint      = "/:collectionName/cover"
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

		{Method: http.MethodGet, Path: collectionStatisticsEndpoint, HandlerFunc: handler.getStatisticsForCollectionWithName},
		{Method: http.MethodPost, Path: collectionCreateEndpoint, HandlerFunc: handler.create},
		{Method: http.MethodGet, Path: collectionAssetsEndpoint, HandlerFunc: handler.getAssetsForCollectionWithName},

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
// @Success 200 {object} []data.Collection
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /collections/list/{offset}/{limit} [get]
func (handler *collectionsHandler) getList(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

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

	collections, err := storage.GetCollectionsWithOffsetLimit(offset, limit)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, collections, "")
}

// @Summary Gets collection by name.
// @Description Retrieves a collection by its name.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionName path string true "collection name"
// @Success 200 {object} data.Collection
// @Failure 404 {object} data.ApiResponse
// @Router /collections/{collectionName} [get]
func (handler *collectionsHandler) get(c *gin.Context) {
	collectionName := c.Param("collectionName")

	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, collection, "")
}

// @Summary Set collection info.
// @Description Sets info for a collection.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionName path string true "collection name"
// @Param updateCollectionRequest body services.UpdateCollectionRequest true "collection info"
// @Success 200 {object} data.Collection
// @Failure 401 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /collections/{collectionName} [post]
func (handler *collectionsHandler) set(c *gin.Context) {
	collectionName := c.Param("collectionName")
	var request services.UpdateCollectionRequest

	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	creator, err := storage.GetAccountById(collection.CreatorID)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if creator.Address != jwtAddress {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	err = services.UpdateCollection(collection, &request)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, collection, "")
}

// @Summary Gets collection statistics.
// @Description Gets statistics for a collection. It will be cached for 10 minutes.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionName path string true "collection name"
// @Success 200 {object} services.CollectionStatistics
// @Failure 404 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /collections/{collectionName}/statistics [post]
func (handler *collectionsHandler) getStatisticsForCollectionWithName(c *gin.Context) {
	collectionName := c.Param("collectionName")

	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	stats, err := services.GetStatisticsForCollection(collection.ID)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, stats, "")
}

// @Summary Creates a collection.
// @Description Creates a collection with given info.
// @Tags collections
// @Accept json
// @Produce json
// @Param createCollectionRequest body services.CreateCollectionRequest true "collection info"
// @Success 200 {object} data.Collection
// @Failure 400 {object} data.ApiResponse
// @Failure 401 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /collections/create [post]
func (handler *collectionsHandler) create(c *gin.Context) {
	var request services.CreateCollectionRequest

	err := c.Bind(&request)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	if jwtAddress != request.UserAddress {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "jwt and request addresses differ")
		return
	}

	collection, err := services.CreateCollection(&request, handler.blockchainCfg.ProxyUrl)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, collection, "")
}

// @Summary Get collection assets.
// @Description Retrieves the assets of a collection. Unsorted.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionName path string true "collection name"
// @Param offset path int true "offset"
// @Param limit path int true "limit"
// @Success 200 {object} []data.Asset
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /collections/{collectionName}/assets/{offset}/{limit} [get]
func (handler *collectionsHandler) getAssetsForCollectionWithName(c *gin.Context) {
	collectionName := c.Param("collectionName")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")
	filters := c.QueryMap("filters")

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

	collection, err := storage.GetCollectionByName(collectionName)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	assets, err := storage.GetAssetsByCollectionIdWithOffsetLimit(collection.ID, offset, limit, filters)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, assets, "")
}

// @Summary Get collection profile image
// @Description Retrieves a collection cover image. It will be sent as base64 encoding (sdt, raw) of its byte representation.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionName path string true "collection name"
// @Success 200 {object} string
// @Failure 404 {object} data.ApiResponse
// @Router /collections/{collectionName}/profile [get]
func (handler *collectionsHandler) getCollectionProfile(c *gin.Context) {
	collectionName := c.Param("collectionName")

	image, err := services.GetCollectionProfileImage(collectionName)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, &image, "")
}

// @Summary Set collection profile image
// @Description Expects base64 encoding (sdt, raw) of the image representation. Returns empty string. Max size of byte array is 1MB.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionName path string true "collection name"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /collections/{collectionName}/profile [post]
func (handler *collectionsHandler) setCollectionProfile(c *gin.Context) {
	var imageBase64 string
	collectionName := c.Param("collectionName")

	err := c.Bind(&imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	err = services.SetCollectionProfileImage(collectionName, &imageBase64, jwtAddress)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, "", "")
}

// @Summary Get collection cover image
// @Description Retrieves a collection cover image. It will be sent as base64 encoding (sdt, raw) of its byte representation.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionName path string true "collection name"
// @Success 200 {object} string
// @Failure 404 {object} data.ApiResponse
// @Router /collections/{collectionName}/cover [get]
func (handler *collectionsHandler) getCollectionCover(c *gin.Context) {
	collectionName := c.Param("collectionName")

	image, err := services.GetCollectionCoverImage(collectionName)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, &image, "")
}

// @Summary Set collection cover image
// @Description Expects base64 encoding (sdt, raw) of the image representation. Returns empty string. Max size of byte array is 1MB.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionName path string true "collection name"
// @Param image body string true "base64 encoded image"
// @Success 200 {object} string
// @Failure 400 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /collections/{collectionName}/cover [post]
func (handler *collectionsHandler) setCollectionCover(c *gin.Context) {
	var imageBase64 string
	collectionName := c.Param("collectionName")

	err := c.Bind(&imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress := c.GetString(middleware.AddressKey)
	err = services.SetCollectionCoverImage(collectionName, &imageBase64, jwtAddress)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, "", "")
}
