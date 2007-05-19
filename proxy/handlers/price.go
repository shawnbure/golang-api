package handlers

import (
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	basePriceEndpoint = "/price"
)

type priceHandler struct {
}

func NewPriceHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &priceHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: "", HandlerFunc: handler.get},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             basePriceEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *priceHandler) get(c *gin.Context) {
	price, err := services.GetEGLDPrice()
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, "could not get price")
		return
	}

	data.JsonResponse(c, http.StatusOK, price, "")
}
