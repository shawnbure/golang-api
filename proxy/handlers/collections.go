package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/services"
	"github.com/erdsea/erdsea-api/stats/collstats"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/gin-gonic/gin"
)

type RankingEntry = collstats.LeaderboardEntry

const (
	baseCollectionsEndpoint    = "/collections"
	collectionByNameEndpoint   = "/:collectionId"
	collectionListEndpoint     = "/list/:offset/:limit"
	collectionCreateEndpoint   = "/create"
	collectionTokensEndpoint   = "/:collectionId/tokens/:offset/:limit"
	collectionProfileEndpoint  = "/:collectionId/profile/"
	collectionCoverEndpoint    = "/:collectionId/cover"
	collectionMintInfoEndpoint = "/:collectionId/mintInfo"
	collectionRankingEndpoint  = "/rankings/:offset/:limit"
)

type collectionsHandler struct {
	blockchainCfg config.BlockchainConfig
}

func NewCollectionsHandler(groupHandler *groupHandler, authCfg config.AuthConfig, blockchainCfg config.BlockchainConfig) {
	handler := &collectionsHandler{blockchainCfg: blockchainCfg}

	endpoints := []EndpointHandler{
		{Method: http.MethodPost, Path: collectionByNameEndpoint, HandlerFunc: handler.set},
		{Method: http.MethodPost, Path: collectionCreateEndpoint, HandlerFunc: handler.create},
		{Method: http.MethodPost, Path: collectionProfileEndpoint, HandlerFunc: handler.setCollectionProfile},
		{Method: http.MethodPost, Path: collectionCoverEndpoint, HandlerFunc: handler.setCollectionCover},
	}
	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseCollectionsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}
	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)

	publicEndpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: collectionListEndpoint, HandlerFunc: handler.getList},
		{Method: http.MethodGet, Path: collectionByNameEndpoint, HandlerFunc: handler.get},
		{Method: http.MethodGet, Path: collectionTokensEndpoint, HandlerFunc: handler.getTokens},
		{Method: http.MethodGet, Path: collectionMintInfoEndpoint, HandlerFunc: handler.getMintInfo},
		{Method: http.MethodGet, Path: collectionRankingEndpoint, HandlerFunc: handler.getCollectionRankings},
	}
	publicEndpointGroupHandler := EndpointGroupHandler{
		Root:             baseCollectionsEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: publicEndpoints,
	}
	groupHandler.AddEndpointGroupHandler(publicEndpointGroupHandler)
}

// @Summary Gets collections.
// @Description Retrieves a list of collections. Sorted by priority.
// @Tags collections
// @Accept json
// @Produce json
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.Collection
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /collections/list/{offset}/{limit} [get]
func (handler *collectionsHandler) getList(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

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

	collections, err := storage.GetCollectionsWithOffsetLimit(int(offset), int(limit))
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
// @Success 200 {object} dtos.ExtendedCollectionDto
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /collections/{collectionId} [get]
func (handler *collectionsHandler) get(c *gin.Context) {
	tokenId := c.Param("collectionId")

	cacheInfo, err := collstats.GetOrAddCollectionCacheInfo(tokenId)
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
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	collectionStats, err := collstats.GetStatisticsForTokenId(tokenId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	extendedDto := dtos.ExtendedCollectionDto{
		Collection:           *collection,
		Statistics:           *collectionStats,
		CreatorWalletAddress: creator.Address,
		CreatorName:          creator.Name,
	}

	dtos.JsonResponse(c, http.StatusOK, extendedDto, "")
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

	err := c.BindJSON(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	cacheInfo, err := collstats.GetOrAddCollectionCacheInfo(tokenId)
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

	err := c.BindJSON(&request)
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
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.Token
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /collections/{collectionId}/tokens/{offset}/{limit} [get]
func (handler *collectionsHandler) getTokens(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")
	filters := c.QueryMap("filters")
	tokenId := c.Param("collectionId")
	sortRules := c.QueryMap("sort")

	acceptedCriteria := map[string]bool{"price_nominal": true, "created_at": true}
	err := testInputSortParams(sortRules, acceptedCriteria)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	cacheInfo, err := collstats.GetOrAddCollectionCacheInfo(tokenId)
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

	tokens, err := storage.GetTokensByCollectionIdWithOffsetLimit(cacheInfo.CollectionId, int(offset), int(limit), filters, sortRules)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, tokens, "")
}

// @Summary Set collection profile image
// @Description Expects base64 std encoding of the image representation. Returns empty string. Max size of byte array is 1MB.
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
	tokenId := c.Param("collectionId")

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	imageBase64 := buf.String()
	cacheInfo, err := collstats.GetOrAddCollectionCacheInfo(tokenId)
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

	link, err := services.SetCollectionProfileImage(tokenId, cacheInfo.CollectionId, &imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, link, "")
}

// @Summary Set collection cover image
// @Description Expects base64 std encoding of the image representation. Returns empty string. Max size of byte array is 1MB.
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
	tokenId := c.Param("collectionId")

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	imageBase64 := buf.String()
	cacheInfo, err := collstats.GetOrAddCollectionCacheInfo(tokenId)
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

	link, err := services.SetCollectionCoverImage(tokenId, cacheInfo.CollectionId, &imageBase64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, link, "")
}

// @Summary Gets mint info about a collection.
// @Description Retrieves max supply and total sold for a collection. Cached for 6 seconds.
// @Tags collections
// @Accept json
// @Produce json
// @Param collectionId path string true "collection id"
// @Success 200 {object} services.MintInfo
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /collections/{collectionId}/mintInfo [get]
func (handler *collectionsHandler) getMintInfo(c *gin.Context) {
	tokenId := c.Param("collectionId")

	cacheInfo, err := collstats.GetOrAddCollectionCacheInfo(tokenId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	collection, err := storage.GetCollectionById(cacheInfo.CollectionId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	if len(collection.ContractAddress) == 0 {
		dtos.JsonResponse(c, http.StatusNotFound, nil, "no contract address")
		return
	}

	mintInfo, err := services.GetMintInfoForContract(collection.ContractAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, "")
		return
	}

	dtos.JsonResponse(c, http.StatusOK, mintInfo, "")
}

// @Summary Get collection rankings
// @Description Acts as a leaderboard. Optionally provide ?sort[criteria]=volumeTraded&sort[mode]=asc
// @Tags collections
// @Accept json
// @Produce json
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} RankingEntry
// @Failure 400 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /collections/rankings/{offset}/{limit} [get]
func (handler *collectionsHandler) getCollectionRankings(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")
	sortParams := c.QueryMap("sort")

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

	if len(sortParams) == 0 {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, "no sorting rules")
		return
	}

	acceptedCriteria := map[string]bool{
		"floorprice":   true,
		"volumetraded": true,
		"itemstotal":   true,
		"ownerstotal":  true,
	}
	err = testInputSortParams(sortParams, acceptedCriteria)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	table := sortParams["criteria"]
	isRev := strings.ToLower(sortParams["mode"]) == "desc"
	entries, err := collstats.GetLeaderboardEntries(table, int(offset), int(limit), isRev)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, entries, "")
}

func testInputSortParams(sortParams map[string]string, acceptedCriteria map[string]bool) error {
	if len(sortParams) == 0 {
		return nil
	}

	if len(sortParams) != 2 {
		return errors.New("bad sorting input len")
	}

	if v, ok := sortParams["mode"]; ok {
		vLower := strings.ToLower(v)
		if vLower != "asc" && vLower != "desc" {
			return errors.New("bad sorting mode")
		}
	} else {
		return errors.New("no sorting mode")
	}

	if v, ok := sortParams["criteria"]; ok {
		vLower := strings.ToLower(v)
		if _, accepted := acceptedCriteria[vLower]; !accepted {
			return errors.New("bad sorting criteria")
		}
	} else {
		return errors.New("no sorting criteria")
	}

	return nil
}
