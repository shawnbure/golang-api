package handlers

import (
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseTokensEndpoint             = "/tokens"
	tokenByTokenIdAndNonceEndpoint = "/:tokenId/:nonce"
)

type tokensHandler struct {
}

func NewTokensHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &tokensHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: tokenByTokenIdAndNonceEndpoint, HandlerFunc: handler.getByTokenIdAndNonce},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseTokensEndpoint,
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
// @Success 200 {object} data.Token
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /assets/{tokenId}/{nonce} [get]
func (handler *tokensHandler) getByTokenIdAndNonce(c *gin.Context) {
	tokenId := c.Param("tokenId")
	nonceString := c.Param("nonce")

	nonce, err := strconv.ParseUint(nonceString, 10, 64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	asset, err := storage.GetTokenByTokenIdAndNonce(tokenId, nonce)
	if err != nil {
		data.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, asset, "")
}
