package handlers

import (
	"encoding/hex"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"net/http"

	"github.com/erdsea/erdsea-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseAuthEndpoint    = "/auth"
	accessAuthEndpoint  = "/access"
	refreshAuthEndpoint = "/refresh"
)

type createTokenRequest struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
	Message   string `json:"message"`
}

type authHandler struct {
	service services.AuthService
}

func RegisterAuthHandler(authService services.AuthService, groupHandler *groupHandler) {
	h := authHandler{
		service: authService,
	}

	endpoints := []EndpointHandler{
		{Method: http.MethodPost, Path: accessAuthEndpoint, HandlerFunc: h.createAccessToken},
		{Method: http.MethodGet, Path: refreshAuthEndpoint, HandlerFunc: h.refreshAccessToken},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseAuthEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (h *authHandler) createAccessToken(c *gin.Context) {
	req := createTokenRequest{}

	var badReqResp = func(c *gin.Context, err string) {
		JsonResponse(c, http.StatusBadRequest, nil, err)
	}

	err := c.Bind(&req)
	if err != nil {
		badReqResp(c, err.Error())
		return
	}

	pk, err := data.NewAddressFromBech32String(req.Address)
	if err != nil {
		badReqResp(c, err.Error())
		return
	}

	sigBytes, err := hex.DecodeString(req.Signature)
	if err != nil {
		badReqResp(c, err.Error())
		return
	}

	msgBytes, err := hex.DecodeString(req.Message)
	if err != nil {
		badReqResp(c, err.Error())
	}

	jwt, refresh, err := h.service.CreateToken(pk.AddressBytes(), sigBytes, msgBytes)
	if err != nil {
		JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	JsonResponse(c, http.StatusOK, gin.H{
		"token":   jwt,
		"refresh": refresh,
	}, "")
}

func (h *authHandler) refreshAccessToken(c *gin.Context) {

}
