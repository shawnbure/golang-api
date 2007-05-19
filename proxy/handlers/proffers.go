package handlers

import (
	"github.com/erdsea/erdsea-api/services"
	"github.com/erdsea/erdsea-api/storage"
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/gin-gonic/gin"
)

const (
	baseProffersEndpoint     = "/proffers"
	proffersForTokenEndpoint = "/:tokenId/:nonce/:offset/:limit"
)

type profferHandler struct {
}

func NewProfferHandler(groupHandler *groupHandler) {
	handler := &profferHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: proffersForTokenEndpoint, HandlerFunc: handler.getProffersForToken},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseProffersEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary Get proffers (offers and bids) for token
// @Description Retrieves proffers for a token (identified by tokenId and nonce)
// @Tags tokens
// @Accept json
// @Produce json
// @Param tokenId path string true "token id"
// @Param nonce path int true "token nonce"
// @Param offset path uint true "offset"
// @Param limit path uint true "limit"
// @Success 200 {object} []entities.Proffer
// @Failure 400 {object} dtos.ApiResponse
// @Failure 404 {object} dtos.ApiResponse
// @Router /proffers/{tokenId}/{nonce}/{offset}/{limit} [get]
func (handler *profferHandler) getProffersForToken(c *gin.Context) {
	tokenId := c.Param("tokenId")
	nonceString := c.Param("nonce")
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")

	nonce, err := strconv.ParseUint(nonceString, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	offset, err := strconv.ParseUint(offsetStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	err = ValidateLimit(limit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	tokenCacheInfo, err := services.GetOrAddTokenCacheInfo(tokenId, nonce)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	proffers, err := storage.GetProffersForTokenWithOffsetLimit(tokenCacheInfo.TokenDbId, int(offset), int(limit))
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, proffers, "")
}
