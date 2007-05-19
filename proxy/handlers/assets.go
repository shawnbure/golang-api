package handlers

import (
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	baseAssetsEndpoint = "/assets"
	getAssetsEndpoint  = "/:id"
)

type assetsHandler struct {
}

func NewAssetsHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	h := &assetsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: getAssetsEndpoint, HandlerFunc: h.getById},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseAssetsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (a *assetsHandler) getById(c *gin.Context) {
	id := c.Param("id")

	data.JsonResponse(c, http.StatusOK, id, "")
}
