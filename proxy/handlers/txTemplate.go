package handlers

import (
	"net/http"
	"strconv"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/formatter"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/stats/collstats"
	"github.com/erdsea/erdsea-api/storage"
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

// @Summary Gets tx-template for NFT list.
// @Description Retrieves tx-template for NFT list. Only account nonce and signature must be added afterwards.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Param tokenId path int true "token id"
// @Param nonce path int true "nonce"
// @Param price path float64 true "price"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/list-nft/{userAddress}/{tokenId}/{nonce}/{price} [get]
func (handler *txTemplateHandler) getListNftTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")
	priceStr := c.Param("price")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template, err := handler.txFormatter.NewListNftTxTemplate(userAddress, tokenId, nonce, price)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Gets tx-template for NFT buy.
// @Description Retrieves tx-template for NFT buy. Only account nonce and signature must be added afterwards.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Param tokenId path int true "token id"
// @Param nonce path int true "nonce"
// @Param price path float64 true "price"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/buy-nft/{userAddress}/{tokenId}/{nonce}/{price} [get]
func (handler *txTemplateHandler) getBuyNftTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")
	priceStr := c.Param("price")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.NewBuyNftTxTemplate(userAddress, tokenId, nonce, priceStr)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Gets tx-template for NFT withdraw.
// @Description Retrieves tx-template for NFT withdraw. Only account nonce and signature must be added afterwards.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Param tokenId path int true "token id"
// @Param nonce path int true "nonce"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/withdraw-nft/{userAddress}/{tokenId}/{nonce} [get]
func (handler *txTemplateHandler) getWithdrawNftTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.NewWithdrawNftTxTemplate(userAddress, tokenId, nonce)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Gets tx-template for mint tokens.
// @Description Retrieves tx-template for mint tokens. Only account nonce and signature must be added afterwards.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Param collectionId path string true "collection id"
// @Param numberOfTokens path int true "number of tokens"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /tx-template/mint-tokens/{userAddress}/{collectionId}/{numberOfTokens} [get]
func (handler *txTemplateHandler) getMintNftTxTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("collectionId")
	numberOfTokensStr := c.Param("numberOfTokens")

	numberOfTokens, err := strconv.ParseUint(numberOfTokensStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	cacheInfo, err := collstats.GetOrAddCollectionCacheInfo(tokenId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	collection, err := storage.GetCollectionById(cacheInfo.CollectionId)
	if err != nil {
		dtos.JsonResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}

	if collection.ContractAddress == "" {
		dtos.JsonResponse(c, http.StatusNotFound, nil, "no contract address")
		return
	}

	pricePerTokenNominal, err := strconv.ParseFloat(collection.MintPricePerTokenString, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, "no contract address")
		return
	}

	template := handler.txFormatter.NewMintNftsTxTemplate(
		userAddress,
		collection.ContractAddress,
		pricePerTokenNominal,
		numberOfTokens,
	)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}
