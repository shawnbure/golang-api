package handlers

import (
	"net/http"

	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseEGLDPriceEndpoint = "/egld_price"
)

type eEGLDPriceHandler struct {
}

func NewPriceHandler(groupHandler *groupHandler) {
	handler := &eEGLDPriceHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: "", HandlerFunc: handler.get},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseEGLDPriceEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Gets EGLD price in dollars.
// @Description Retrieves EGLD price in dollars. Price taken from Binance. Cached for 15 minutes.
// @Tags egld_price
// @Accept json
// @Produce json
// @Success 200 {object} float64
// @Failure 500 {object} dtos.ApiResponse
// @Router /egld_price [get]
func (handler *eEGLDPriceHandler) get(c *gin.Context) {
	price, err := services.GetEGLDPrice()
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, "could not get price")
		return
	}

	dtos.JsonResponse(c, http.StatusOK, price, "")
}
