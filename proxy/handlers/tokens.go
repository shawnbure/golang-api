package handlers

import (
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/services"
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

// @Summary Get token by id and nonce
// @Description Retrieves a token by tokenId and nonce
// @Tags tokens
// @Accept json
// @Produce json
// @Param tokenId path string true "token id"
// @Param nonce path int true "token nonce"
// @Success 200 {object} data.Token
// @Failure 400 {object} data.ApiResponse
// @Failure 404 {object} data.ApiResponse
// @Router /tokens/{tokenId}/{nonce} [get]
func (handler *tokensHandler) getByTokenIdAndNonce(c *gin.Context) {
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
