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

type refreshTokenRequest struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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

	err := c.Bind(&req)
	if err != nil {
		h.badReqResp(c, err.Error())
		return
	}

	pk, err := data.NewAddressFromBech32String(req.Address)
	if err != nil {
		h.badReqResp(c, err.Error())
		return
	}

	sigBytes, err := hex.DecodeString(req.Signature)
	if err != nil {
		h.badReqResp(c, err.Error())
		return
	}

	msgBytes, err := hex.DecodeString(req.Message)
	if err != nil {
		h.badReqResp(c, err.Error())
	}

	jwt, refresh, err := h.service.CreateToken(pk.AddressBytes(), sigBytes, msgBytes)
	if err != nil {
		JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	JsonResponse(c, http.StatusOK, gin.H{
		"accessToken":  jwt,
		"refreshToken": refresh,
	}, "")
}

func (h *authHandler) refreshAccessToken(c *gin.Context) {
	req := refreshTokenRequest{}

	err := c.Bind(&req)
	if err != nil {
		h.badReqResp(c, err.Error())
		return
	}

	jwt, refresh, err := h.service.RefreshToken(req.AccessToken, req.RefreshToken)
	if err != nil {
		JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	JsonResponse(c, http.StatusOK, gin.H{
		"accessToken":  jwt,
		"refreshToken": refresh,
	}, "")
}

func (h *authHandler) badReqResp(c *gin.Context, err string) {
	JsonResponse(c, http.StatusBadRequest, nil, err)
}
