package handlers

import (
	"net/http"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseDepositsEndpoint      = "/deposits"
	depositForAddressEndpoint = "/:userAddress"
)

type depositsHandler struct {
	cfg config.BlockchainConfig
}

func NewDepositsHandler(groupHandler *groupHandler, cfg config.BlockchainConfig) {
	handler := &depositsHandler{cfg: cfg}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: depositForAddressEndpoint, HandlerFunc: handler.getDepositForAddress},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseDepositsEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Gets the deposit (EGLD) located in the marketplace for an address.
// @Description Retrieves deposit amount for an address.
// @Tags deposits
// @Accept json
// @Produce json
// @Param userAddress path string true "userAddress"
// @Success 200 {object} float64
// @Failure 400 {object} dtos.ApiResponse
// @Router /deposits/{userAddress} [get]
func (handler *depositsHandler) getDepositForAddress(c *gin.Context) {
	userAddress := c.Param("userAddress")

	_, err := services.GetOrAddAccountCacheInfo(userAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	deposit, err := services.GetDeposit(handler.cfg.MarketplaceAddress, userAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, deposit, "")
}
