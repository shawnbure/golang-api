package handlers

import (
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/formatter"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/gin-gonic/gin"
)

const (
	baseFormatEndpoint        = "/tx-template"
	listNftFormatEndpoint     = "/list-nft/:userAddress/:tokenId/:nonce/:price"
	buyNftFormatEndpoint      = "/buy-nft/:userAddress/:tokenId/:nonce/:price"
	withdrawNftFormatEndpoint = "/withdraw-nft/:userAddress/:tokenId/:nonce"
)

type txTemplateHandler struct {
	txFormatter formatter.TxFormatter
}

func NewTxTemplateHandler(groupHandler *groupHandler, authCfg config.AuthConfig, blockchainConfig config.BlockchainConfig) {
	handler := &txTemplateHandler{
		txFormatter: formatter.NewTxFormatter(blockchainConfig),
	}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: listNftFormatEndpoint, HandlerFunc: handler.getListNftTemplate},
		{Method: http.MethodGet, Path: buyNftFormatEndpoint, HandlerFunc: handler.getBuyNftTemplate},
		{Method: http.MethodGet, Path: withdrawNftFormatEndpoint, HandlerFunc: handler.getWithdrawNftTemplate},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseFormatEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *txTemplateHandler) getListNftTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")
	priceStr := c.Param("price")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template, err := handler.txFormatter.NewListNftTxTemplate(userAddress, tokenId, nonce, price)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, template, "")
}

func (handler *txTemplateHandler) getBuyNftTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")
	priceStr := c.Param("price")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.NewBuyNftTxTemplate(userAddress, tokenId, nonce, price)

	data.JsonResponse(c, http.StatusOK, template, "")
}

func (handler *txTemplateHandler) getWithdrawNftTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.NewWithdrawNftTxTemplate(userAddress, tokenId, nonce)
	data.JsonResponse(c, http.StatusOK, template, "")
}
