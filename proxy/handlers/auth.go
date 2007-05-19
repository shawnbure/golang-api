package handlers

import (
	"encoding/hex"
	"net/http"

	erdgoData "github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/erdsea/erdsea-api/data"
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

type tokenPayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type authHandler struct {
	service services.AuthService
}

func NewAuthHandler(groupHandler *groupHandler, authService services.AuthService) {
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

// @Summary access credentials
// @Description creates an access credentials
// @Tags auth
// @Accept  json
// @Produce  json
// @Param tokenRequest body createTokenRequest true "create credentials request"
// @Success 200 {object} tokenPayload
// @Failure 400 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /auth/access [post]
func (h *authHandler) createAccessToken(c *gin.Context) {
	req := createTokenRequest{}

	err := c.Bind(&req)
	if err != nil {
		h.badReqResp(c, err.Error())
		return
	}

	pk, err := erdgoData.NewAddressFromBech32String(req.Address)
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
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, tokenPayload{
		AccessToken:  jwt,
		RefreshToken: refresh,
	}, "")
}

// @Summary refresh credentials
// @Description refreshes the access credentials
// @Tags auth
// @Accept  json
// @Produce  json
// @Param refreshRequest body tokenPayload true "refresh credentials request"
// @Success 200 {object} tokenPayload
// @Failure 400 {object} data.ApiResponse
// @Failure 500 {object} data.ApiResponse
// @Router /auth/refresh [post]
func (h *authHandler) refreshAccessToken(c *gin.Context) {
	req := tokenPayload{}

	err := c.Bind(&req)
	if err != nil {
		h.badReqResp(c, err.Error())
		return
	}

	jwt, refresh, err := h.service.RefreshToken(req.AccessToken, req.RefreshToken)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, tokenPayload{
		AccessToken:  jwt,
		RefreshToken: refresh,
	}, "")
}

func (h *authHandler) badReqResp(c *gin.Context, err string) {
	data.JsonResponse(c, http.StatusBadRequest, nil, err)
}
