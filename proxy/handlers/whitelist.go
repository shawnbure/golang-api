package handlers

import (
	"net/http"
	"strconv"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/proxy/middleware"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseWhitelistEndpoint    = "/whitelists"
	whitelistByAddress       = "/:address"
	whitelistByAddressAmount = "/:address/:amount"
	//whitelistBuyLimitByAddresses = "buy-limit/:contract-address/:user-address"
)

type whitelistHandler struct {
	blockchainCfg config.BlockchainConfig
}

func NewWhitelistHandler(groupHandler *groupHandler, authCfg config.AuthConfig, blockchainCfg config.BlockchainConfig) {
	handler := &whitelistHandler{blockchainCfg: blockchainCfg}

	endpoints := []EndpointHandler{
		{Method: http.MethodPost, Path: whitelistByAddress, HandlerFunc: handler.GetWhitelistByAddress},
		{Method: http.MethodPost, Path: whitelistByAddressAmount, HandlerFunc: handler.UpdateWhitelistAmountByAddress},
		/*{Method: http.MethodPost, Path: whitelistLimitByAddresses, HandlerFunc: handler.GetWhitelistLimitByAddresses},*/
	}
	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseCollectionsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}
	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
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
func (handler *whitelistHandler) UpdateWhitelistAmountByAddress(c *gin.Context) {
	var request services.SetWhitelistRequest
	address := c.Param("address")
	strAmount := c.Param("amount")

	err := c.BindJSON(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	whitelist, err := storage.GetWhitelistByAddress(address)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	//parse the string amount to Uint amount
	amount, err := strconv.ParseUint(strAmount, 10, 64)

	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	//err = services.UpdateWhitelistAmountByAddress(amount, address)
	err = storage.UpdateWhitelistAmountByAddress(amount, address)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, whitelist, "")
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
func (handler *whitelistHandler) GetWhitelistByAddress(c *gin.Context) {
	address := c.Param("address")

	whitelist, err := services.GetWhitelist(address) //storage.GetWhitelistByAddress(address)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, whitelist, "")
}
