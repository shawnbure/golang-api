package handlers

import (
	"github.com/erdsea/erdsea-api/services"
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/gin-gonic/gin"
)

const (
	baseAssetsEndpoint             = "/assets"
	assetByTokenIdAndNonceEndpoint = "/:tokenId/:nonce"
)

type assetsHandler struct {
}

func NewAssetsHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &assetsHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: assetByTokenIdAndNonceEndpoint, HandlerFunc: handler.getByTokenIdAndNonce},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseAssetsEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Get asset by token by id and nonce
// @Description Retrieves an asset by tokenId and nonce
// @Tags assets
// @Accept json
// @Produce json
// @Param tokenId path string true "token id"
// @Param nonce path int true "token nonce"
// @Success 200 {object} data.Asset
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /assets/{tokenId}/{nonce} [get]
func (handler *assetsHandler) getByTokenIdAndNonce(c *gin.Context) {
	tokenId := c.Param("tokenId")
	nonceString := c.Param("nonce")

	nonce, err := strconv.ParseUint(nonceString, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	tokenDto, err := services.GetExtendedTokenData(tokenId, nonce)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, tokenDto, "")
}
