package handlers

import (
	"net/http"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseEGLDPriceEndpoint = "/egld_price"
)

type eEGLDPriceHandler struct {
}

func NewPriceHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &eEGLDPriceHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: "", HandlerFunc: handler.get},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseEGLDPriceEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
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
// @Failure 500 {object} data.ApiResponse
// @Router /egld_price [get]
func (handler *eEGLDPriceHandler) get(c *gin.Context) {
	price, err := services.GetEGLDPrice()
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, "could not get price")
		return
	}

	data.JsonResponse(c, http.StatusOK, price, "")
}
