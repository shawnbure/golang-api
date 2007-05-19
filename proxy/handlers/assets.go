package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	baseAssetsEndpoint = "/assets"
	getAssetsEndpoint  = "/:id"
)

type assetsHandler struct {
}

func NewAssetsHandler(groupHandler *groupHandler) {
	h := &assetsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: getAssetsEndpoint, HandlerFunc: h.getById},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseAssetsEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (a *assetsHandler) getById(c *gin.Context) {
	id := c.Param("id")

	JsonResponse(c, http.StatusOK, id, "")
}
