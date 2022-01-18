package handlers

import (
	"net/http"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseRoyaltiesEndpoint                         = "/royalties"
	royaltiesForAddressEndpoint                   = "/:userAddress/amount"
	lastWithdrawalEpochForAddressEndpoint         = "/:userAddress/last"
	remainingEpochUntilWithdrawForAddressEndpoint = "/:userAddress/remaining"
)

type royaltiesHandler struct {
	cfg config.BlockchainConfig
}

func NewRoyaltiesHandler(groupHandler *groupHandler, cfg config.BlockchainConfig) {
	handler := &royaltiesHandler{cfg: cfg}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: royaltiesForAddressEndpoint, HandlerFunc: handler.getRoyaltiesForAddress},
		{Method: http.MethodGet, Path: lastWithdrawalEpochForAddressEndpoint, HandlerFunc: handler.getLastWithdrawalEpochForAddress},
		{Method: http.MethodGet, Path: remainingEpochUntilWithdrawForAddressEndpoint, HandlerFunc: handler.getRemainingEpochsUntilWithdrawForAddress},
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
// @Router /royalties/{userAddress}/amount [get]
func (handler *royaltiesHandler) getRoyaltiesForAddress(c *gin.Context) {
	userAddress := c.Param("userAddress")

	deposit, err := services.GetCreatorRoyalties(handler.cfg.MarketplaceAddress, userAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, deposit, "")
}

// @Summary Gets last withdrawal epoch (EGLD) for an address.
// @Description Gets last withdrawal epoch for a creator.
// @Tags royalties
// @Accept json
// @Produce json
// @Param userAddress path string true "userAddress"
// @Success 200 {object} int64
// @Failure 400 {object} dtos.ApiResponse
// @Router /royalties/{userAddress}/last [get]
func (handler *royaltiesHandler) getLastWithdrawalEpochForAddress(c *gin.Context) {
	userAddress := c.Param("userAddress")

	deposit, err := services.GetCreatorLastWithdrawalEpoch(handler.cfg.MarketplaceAddress, userAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, deposit, "")
}

// @Summary Gets remaining epochs until withdraw royalties for an address.
// @Description Gets remaining epochs until withdrawal epoch for a creator.
// @Tags royalties
// @Accept json
// @Produce json
// @Param userAddress path string true "userAddress"
// @Success 200 {object} int64
// @Failure 400 {object} dtos.ApiResponse
// @Router /royalties/{userAddress}/remaining [get]
func (handler *royaltiesHandler) getRemainingEpochsUntilWithdrawForAddress(c *gin.Context) {
	userAddress := c.Param("userAddress")

	epochs, err := services.GetCreatorRemainingEpochsUntilWithdraw(handler.cfg.MarketplaceAddress, userAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, epochs, "")
}
