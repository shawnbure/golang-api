package handlers

import (
	"net/http"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseRoyaltiesEndpoint                 = "/royalties"
	royaltiesForAddressEndpoint           = "/:userAddress"
	lastWithdrawalEpochForAddressEndpoint = "/last/:userAddress"
)

type royaltiesHandler struct {
	cfg config.BlockchainConfig
}

func NewRoyaltiesHandler(groupHandler *groupHandler, cfg config.BlockchainConfig) {
	handler := &royaltiesHandler{cfg: cfg}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: royaltiesForAddressEndpoint, HandlerFunc: handler.getRoyaltiesForAddress},
		{Method: http.MethodGet, Path: lastWithdrawalEpochForAddressEndpoint, HandlerFunc: handler.getLastWithdrawalEpochForAddress},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseRoyaltiesEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Gets the royalties (EGLD) located in the marketplace for an address.
// @Description Retrieves royalties amount for an address.
// @Tags royalties
// @Accept json
// @Produce json
// @Param userAddress path string true "userAddress"
// @Success 200 {object} float64
// @Failure 400 {object} dtos.ApiResponse
// @Router /royalties/{userAddress} [get]
func (handler *royaltiesHandler) getRoyaltiesForAddress(c *gin.Context) {
	userAddress := c.Param("userAddress")

	_, err := services.GetOrAddAccountCacheInfo(userAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	deposit, err := services.GetCreatorRoyalties(handler.cfg.MarketplaceAddress, userAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, deposit, "")
}

// @Summary Gets last withdrawal epoch (EGLD) for an address.
// @Description Gets last withdrawal epoch for a creator. Next withdraw needs to be calculated as (current_epoch - this_epoch) %30
// @Tags royalties
// @Accept json
// @Produce json
// @Param userAddress path string true "userAddress"
// @Success 200 {object} float64
// @Failure 400 {object} dtos.ApiResponse
// @Router /royalties/last/{userAddress} [get]
func (handler *royaltiesHandler) getLastWithdrawalEpochForAddress(c *gin.Context) {
	userAddress := c.Param("userAddress")

	_, err := services.GetOrAddAccountCacheInfo(userAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	deposit, err := services.GetCreatorLastWithdrawalEpoch(handler.cfg.MarketplaceAddress, userAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, deposit, "")
}
