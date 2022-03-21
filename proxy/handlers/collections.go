package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/proxy/middleware"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/ENFT-DAO/youbei-api/stats/collstats"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/gin-gonic/gin"
)

type RankingEntry = collstats.LeaderboardEntry

const (
	baseCollectionsEndpoint      = "/collections"
	collectionByNameEndpoint     = "/:collectionId"
	collectionListEndpoint       = "/list/:offset/:limit"
	collectionCreateEndpoint     = "/create"
	collectionTokensEndpoint     = "/:collectionId/tokens/:offset/:limit"
	collectionProfileEndpoint    = "/:collectionId/profile"
	collectionCoverEndpoint      = "/:collectionId/cover"
	collectionMintInfoEndpoint   = "/:collectionId/mintInfo"
	collectionRankingEndpoint    = "/rankings/:offset/:limit"
	collectionAllEndpoint        = "/all"
	collectionVerifiedEndpoint   = "/verified/:limit"
	collectionNoteworthyEndpoint = "/noteworthy/:limit"
	collectionTrendingEndpoint   = "/trending/:limit"
)

type CollectionTokensQueryBody struct {
	Filters   map[string]string `json:"filters"`
	SortRules map[string]string `json:"sortRules"`
}

type CollectionRankingQueryBody struct {
	SortRules map[string]string `json:"sortRules"`
}

type CollectionListQueryBody struct {
	Flags []string `json:"flags"`
}

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
		{Method: http.MethodPost, Path: collectionListEndpoint, HandlerFunc: handler.getList},
		{Method: http.MethodGet, Path: collectionByNameEndpoint, HandlerFunc: handler.get},
		{Method: http.MethodPost, Path: collectionTokensEndpoint, HandlerFunc: handler.getTokens},
		{Method: http.MethodGet, Path: collectionMintInfoEndpoint, HandlerFunc: handler.getMintInfo},
		{Method: http.MethodPost, Path: collectionRankingEndpoint, HandlerFunc: handler.getCollectionRankings},
		{Method: http.MethodPost, Path: collectionAllEndpoint, HandlerFunc: handler.getAll},
		{Method: http.MethodGet, Path: collectionVerifiedEndpoint, HandlerFunc: handler.getCollectionVerified},
		{Method: http.MethodGet, Path: collectionNoteworthyEndpoint, HandlerFunc: handler.getCollectionNoteworthy},
		{Method: http.MethodGet, Path: collectionTrendingEndpoint, HandlerFunc: handler.getCollectionTrending},
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
// @Param query body CollectionListQueryBody true "flag array"
// @Success 200 {object} []entities.Collection
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /collections/list/{offset}/{limit} [post]
func (handler *collectionsHandler) getList(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	var queries CollectionListQueryBody
	err := c.BindJSON(&queries)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	flags := queries.Flags

	err = services.CheckValidFlags(flags)
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

	collections, err := storage.GetCollectionsWithOffsetLimit(int(offset), int(limit), flags)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, collections, "")
}

func (handler *collectionsHandler) getAll(c *gin.Context) {

	collections, err := services.GetAllCollections()
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, collections, "")
}

func (handler *collectionsHandler) getCollectionVerified(c *gin.Context) {

	limitStr := c.Param("limit")

	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	collections, err := services.GetCollectionsVerified(int(limit))
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, collections, "")
}

func (handler *collectionsHandler) getCollectionNoteworthy(c *gin.Context) {

	limitStr := c.Param("limit")
	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	collections, err := services.GetCollectionsNoteworthy(int(limit))
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, collections, "")
}

func (handler *collectionsHandler) getCollectionTrending(c *gin.Context) {

	limitStr := c.Param("limit")
	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	collections, err := services.GetCollectionsTrending(int(limit))
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
// @Param collectionId path string true "colle`ction id"
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
// @Param query body CollectionTokensQueryBody true "filters and sort rules"
// @Success 200 {object} []entities.Token
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /collections/{collectionId}/tokens/{offset}/{limit} [post]
func (handler *collectionsHandler) getTokens(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")
	tokenId := c.Param("collectionId")

	var queries CollectionTokensQueryBody
	err := c.BindJSON(&queries)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	sortRules := queries.SortRules
	filters := queries.Filters

	acceptedCriteria := map[string]bool{"price_nominal": true, "created_at": true}
	err = testInputSortParams(sortRules, acceptedCriteria)
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
// @Description Acts as a leaderboard
// @Tags collections
// @Accept json
// @Produce json
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Param query body CollectionRankingQueryBody true "sort rules"
// @Success 200 {object} RankingEntry
// @Failure 400 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /collections/rankings/{offset}/{limit} [post]
func (handler *collectionsHandler) getCollectionRankings(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	var queries CollectionRankingQueryBody
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
		sortRules["criteria"] = collstats.VolumeTraded
		sortRules["mode"] = "desc"
	}

	acceptedCriteria := map[string]bool{
		collstats.FloorPrice:   true,
		collstats.VolumeTraded: true,
		collstats.ItemsTotal:   true,
		collstats.OwnersTotal:  true,
	}
	err = testInputSortParams(sortRules, acceptedCriteria)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	table := sortRules["criteria"]
	isRev := strings.ToLower(sortRules["mode"]) == "desc"
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
		if _, accepted := acceptedCriteria[v]; !accepted {
			return errors.New("bad sorting criteria")
		}
	} else {
		return errors.New("no sorting criteria")
	}

	return nil
}
