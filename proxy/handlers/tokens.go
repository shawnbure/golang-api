package handlers

import (
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseTokensEndpoint             = "/tokens"
	tokenByTokenIdAndNonceEndpoint = "/:tokenId/:nonce"
	availableTokensEndpoint        = "/available"
)

type AvailableTokensRequest struct {
	Tokens []string `json:"tokens"`
}

type AvailableToken struct {
	Collection struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"collection"`
	Token struct {
		Id        string `json:"id"`
		Nonce     uint64 `json:"nonce"`
		Name      string `json:"name"`
		Available bool   `json:"available"`
	}
}

type AvailableTokensResponse struct {
	Tokens map[string]AvailableToken `json:"tokens"`
}

type tokensHandler struct {
}

func NewTokensHandler(groupHandler *groupHandler) {
	handler := &tokensHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: tokenByTokenIdAndNonceEndpoint, HandlerFunc: handler.getByTokenIdAndNonce},
		{Method: http.MethodGet, Path: availableTokensEndpoint, HandlerFunc: handler.getAvailableTokens},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseTokensEndpoint,
		Middlewares:      []gin.HandlerFunc{},
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
// @Success 200 {object} entities.Token
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
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

// @Summary Get available tokens
// @Description Get available tokens and some collection info
// @Tags tokens
// @Accept json
// @Produce json
// @Param availableTokensRequest body AvailableTokensRequest true "request"
// @Success 200 {object} AvailableTokensResponse
// @Failure 400 {object} dtos.ApiResponse
// @Router /tokens/available [get]
func (handler *tokensHandler) getAvailableTokens(c *gin.Context) {
	var request AvailableTokensRequest

	err := c.Bind(&request)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	response := services.GetAvailableTokens(request)
	dtos.JsonResponse(c, http.StatusOK, response, "")
}
