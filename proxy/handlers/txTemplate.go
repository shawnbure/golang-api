package handlers

import (
	"fmt"
	"log"
	"os"

	"net/http"
	"strconv"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/formatter"
	"github.com/ENFT-DAO/youbei-api/stats/collstats"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/gin-gonic/gin"
)

const (
	baseFormatEndpoint                         = "/tx-template"
	listNftFormatEndpoint                      = "/list-nft/:userAddress/:tokenId/:nonce/:price"
	buyNftFormatEndpoint                       = "/buy-nft/:userAddress/:tokenId/:nonce/:price"
	withdrawNftFormatEndpoint                  = "/withdraw-nft/:userAddress/:tokenId/:nonce"
	stakeNftFormatEndpoint                     = "/stake/:userAddress/:collectionId/:nonce"
	unstakeNftFormatEndpoint                   = "/unstake/:userAddress/:collectionId/:nonce"
	stakeCollectionFormatEndpoint              = "/stakecol/:userAddress/:collectionId"
	unstakeCollectionFormatEndpoint            = "/unstakecol/:userAddress/:collectionId"
	makeOfferFormatEndpoint                    = "/make-offer/:userAddress/:tokenId/:nonce/:amount/:expire"
	acceptOfferFormatEndpoint                  = "/accept-offer/:userAddress/:tokenId/:nonce/:offerorAddress/:amount"
	cancelOfferFormatEndpoint                  = "/cancel-offer/:userAddress/:tokenId/:nonce/:amount"
	startAuctionFormatEndpoint                 = "/start-auction/:userAddress/:tokenId/:nonce/:minBid/:startTime/:deadline"
	placeBidFormatEndpoint                     = "/place-bid/:userAddress/:tokenId/:nonce/:payment/:bidAmount"
	endAuctionFormatEndpoint                   = "/end-auction/:userAddress/:tokenId/:nonce"
	depositFormatEndpoint                      = "/deposit/:userAddress/:amount"
	mintTokensFormatEndpoint                   = "/mint-tokens/:userAddress/:tokenId/:numberOfTokens"
	withdrawFormatEndpoint                     = "/withdraw/:userAddress/:amount"
	withdrawCreatorRoyaltiesFormatEndpoint     = "/withdraw-creator-royalties/:userAddress"
	redeemStakingRewardsFormatEndpoint         = "/redeem-staking-rewards/:userAddress"
	issueNFTFormatEndpoint                     = "/issue-nft/:userAddress/:tokenName/:tokenTicker"
	deployNFTTemplateFormatEndpoint            = "/deploy-template/:userAddress/:tokenId/:royalties/:tokenNameBase/:imageExt/:price/:maxSupply/:saleStart"
	changeOwnerFormatEndpoint                  = "/change-owner/:userAddress/:contractAddress"
	setSpecialRolesFormatEndpoint              = "/set-roles/:userAddress/:tokenId/:contractAddress"
	withdrawFromMinterFormatEndpoint           = "/withdraw-minter/:userAddress/:contractAddress"
	requestWithdrawThroughMinterFormatEndpoint = "/request-withdraw/:userAddress/:contractAddress"
	updateSaleStartFormatEndpoint              = "/update-sale-start/:userAddress/:contractAddress/:saleStart"
	updateBuyerWhiteListCheckFormatEndpoint    = "/update-buyer-whitelist-check/:userAddress/:contractAddress/:whiteListCheck"
	getBuyerWhiteListCheckFormatEndpoint       = "/get-buyer-whitelist-check/:userAddress/:contractAddress"
	payCheckoutFromDepositFormatEndpoint	   = "/pay-checkout/:userAddress/:amount"
)

type txTemplateHandler struct {
	txFormatter formatter.TxFormatter
	logErr      *log.Logger
}

func NewTxTemplateHandler(groupHandler *groupHandler, blockchainConfig config.BlockchainConfig) {
	l := log.New(os.Stderr, "", log.LUTC|log.LstdFlags|log.Lshortfile)
	handler := &txTemplateHandler{
		txFormatter: formatter.NewTxFormatter(blockchainConfig),
		logErr:      l,
	}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: listNftFormatEndpoint, HandlerFunc: handler.getListNftTemplate},
		{Method: http.MethodGet, Path: buyNftFormatEndpoint, HandlerFunc: handler.getBuyNftTemplate},
		{Method: http.MethodGet, Path: withdrawNftFormatEndpoint, HandlerFunc: handler.getWithdrawNftTemplate},
		{Method: http.MethodGet, Path: stakeNftFormatEndpoint, HandlerFunc: handler.getStakeNftTemplate},
		{Method: http.MethodGet, Path: unstakeNftFormatEndpoint, HandlerFunc: handler.getUnstakeNftTemplate},
		{Method: http.MethodGet, Path: stakeCollectionFormatEndpoint, HandlerFunc: handler.getStakeCollectionTemplate},
		{Method: http.MethodGet, Path: unstakeCollectionFormatEndpoint, HandlerFunc: handler.getUnstakeCollectionTemplate},
		{Method: http.MethodGet, Path: makeOfferFormatEndpoint, HandlerFunc: handler.getMakeOfferTemplate},
		{Method: http.MethodGet, Path: acceptOfferFormatEndpoint, HandlerFunc: handler.getAcceptOfferTemplate},
		{Method: http.MethodGet, Path: cancelOfferFormatEndpoint, HandlerFunc: handler.getCancelOfferTemplate},
		{Method: http.MethodGet, Path: startAuctionFormatEndpoint, HandlerFunc: handler.getStartAuctionTemplate},
		{Method: http.MethodGet, Path: placeBidFormatEndpoint, HandlerFunc: handler.getPlaceBidTemplate},
		{Method: http.MethodGet, Path: endAuctionFormatEndpoint, HandlerFunc: handler.getEndAuctionTemplate},
		{Method: http.MethodGet, Path: depositFormatEndpoint, HandlerFunc: handler.getDepositTemplate},
		{Method: http.MethodGet, Path: withdrawFormatEndpoint, HandlerFunc: handler.getWithdrawTemplate},
		{Method: http.MethodGet, Path: mintTokensFormatEndpoint, HandlerFunc: handler.getMintNftTxTemplate},
		{Method: http.MethodGet, Path: withdrawCreatorRoyaltiesFormatEndpoint, HandlerFunc: handler.getWithdrawCreatorRoyalties},
		{Method: http.MethodGet, Path: redeemStakingRewardsFormatEndpoint, HandlerFunc: handler.getRedeemStakingRewards},
		{Method: http.MethodGet, Path: issueNFTFormatEndpoint, HandlerFunc: handler.getIssueNFT},
		{Method: http.MethodGet, Path: deployNFTTemplateFormatEndpoint, HandlerFunc: handler.getDeployNFTTemplate},
		{Method: http.MethodGet, Path: changeOwnerFormatEndpoint, HandlerFunc: handler.getChangeOwner},
		{Method: http.MethodGet, Path: setSpecialRolesFormatEndpoint, HandlerFunc: handler.getSetSpecialRoles},
		{Method: http.MethodGet, Path: withdrawFromMinterFormatEndpoint, HandlerFunc: handler.withdrawFromMinter},
		{Method: http.MethodGet, Path: requestWithdrawThroughMinterFormatEndpoint, HandlerFunc: handler.requestWithdrawThroughMinter},
		{Method: http.MethodGet, Path: updateSaleStartFormatEndpoint, HandlerFunc: handler.updateSaleStart},
		{Method: http.MethodGet, Path: updateBuyerWhiteListCheckFormatEndpoint, HandlerFunc: handler.updateBuyerWhiteListCheck},
		{Method: http.MethodGet, Path: getBuyerWhiteListCheckFormatEndpoint, HandlerFunc: handler.getBuyerWhiteListCheck},
		{Method: http.MethodGet, Path: payCheckoutFromDepositFormatEndpoint, HandlerFunc: handler.payCheckoutFromDeposit},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseFormatEndpoint,
		Middlewares:      []gin.HandlerFunc{},
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
// @Failure 500 {object} dtos.ApiResponse
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
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
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

	fmt.Print(priceStr)

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	collection, errCollection := storage.GetCollectionByTokenId(tokenId)

	if errCollection != nil { //check if there is an error in getting the collection by token id

		dtos.JsonResponse(c, http.StatusBadRequest, nil, errCollection.Error())
		return
	} else if collection.Type == entities.Collection_type_whitelisted { //collection type is "whitelisted"

		// //Henry - add the check for whitelist here
		// whitelist, errWhitelist := storage.GetWhitelistByAddress(userAddress)

		// //TODO: check if it exist for whitelist, if not, enter in "Not part of the whitelist"
		// if errWhitelist == gorm.ErrRecordNotFound {

		// 	//throw error message "Not Part of the Whitelist"
		// 	dtos.JsonResponse(c, http.StatusBadRequest, nil, "Sorry, you are not part of the whitelist.")
		// 	return
		// } else if whitelist.Amount == 0 {

		// 	//throw an error : "You already bought the allocated amount for the whitelist"
		// 	dtos.JsonResponse(c, http.StatusBadRequest, nil, "You already bought the allocated amount for the whitelist.")
		// 	return
		// } else {

		// 	//deduct the amount by 1
		// 	newAmount := whitelist.Amount - 1

		// 	//update it
		// 	storage.UpdateWhitelistAmountByAddress(uint64(newAmount), userAddress)
		// }

	}

	//TODO: Address signature if needed for buyying. Added place holder for now.
	template := handler.txFormatter.NewBuyNftTxTemplate(userAddress, tokenId, nonce, []byte(""), priceStr)
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

// @Summary Make offer for an NFT - tx template.
// @Description Retrieves tx-template for make offer transaction.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Param tokenId path int true "token id"
// @Param nonce path int true "nonce"
// @Param amount path float64 true "amount"
// @Param expire path int true "nonce"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/make-offer/{userAddress}/{tokenId}/{nonce}/{amount}/{expire} [get]
func (handler *txTemplateHandler) getMakeOfferTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")
	amountStr := c.Param("amount")
	expireStr := c.Param("expire")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	expire, err := strconv.ParseUint(expireStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.MakeOfferTxTemplate(userAddress, tokenId, nonce, amount, expire)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Accepts offer for an NFT - tx template.
// @Description Retrieves tx-template for accept offer transaction.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Param tokenId path int true "token id"
// @Param nonce path int true "nonce"
// @Param offerorAddress path string true "offerorAddress"
// @Param amount path float64 true "amount"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /tx-template/accept-offer/{userAddress}/{tokenId}/{nonce}/{offerorAddress}/{amount} [get]
func (handler *txTemplateHandler) getAcceptOfferTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")
	offerorAddress := c.Param("offerorAddress")
	amountStr := c.Param("amount")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template, err := handler.txFormatter.AcceptOfferTxTemplate(userAddress, tokenId, nonce, offerorAddress, amount)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Cancels offer for an NFT - tx template.
// @Description Retrieves tx-template for cancel offer transaction.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Param tokenId path int true "token id"
// @Param nonce path int true "nonce"
// @Param amount path float64 true "amount"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/cancel-offer/{userAddress}/{tokenId}/{nonce}/{amount} [get]
func (handler *txTemplateHandler) getCancelOfferTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")
	amountStr := c.Param("amount")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.CancelOfferTxTemplate(userAddress, tokenId, nonce, amount)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Start auction for an NFT - tx template.
// @Description Retrieves tx-template for start auction transaction.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Param tokenId path int true "token id"
// @Param nonce path int true "nonce"
// @Param minBid path float64 true "minBid"
// @Param startTime path int true "nonce"
// @Param deadline path int true "nonce"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /tx-template/start-auction/{userAddress}/{tokenId}/{nonce}/{minBid}/{startTime}/{deadline} [get]
func (handler *txTemplateHandler) getStartAuctionTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")
	minBidStr := c.Param("minBid")
	startTimeStr := c.Param("startTime")
	deadlineStr := c.Param("deadline")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	minBid, err := strconv.ParseFloat(minBidStr, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	startTime, err := strconv.ParseUint(startTimeStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	deadline, err := strconv.ParseUint(deadlineStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template, err := handler.txFormatter.StartAuctionTxTemplate(userAddress, tokenId, nonce, minBid, startTime, deadline)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Start auction for an NFT - tx template.
// @Description Retrieves tx-template for place bid transaction.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Param tokenId path int true "token id"
// @Param nonce path int true "nonce"
// @Param payment path float64 true "payment"
// @Param bidAmount path float64 true "bidAmount"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/place-bid/{userAddress}/{tokenId}/{nonce}/{payment}/{bidAmount} [get]
func (handler *txTemplateHandler) getPlaceBidTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")
	paymentStr := c.Param("payment")
	bidAmountStr := c.Param("bidAmount")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	bidAmount, err := strconv.ParseFloat(bidAmountStr, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.PlaceBidTxTemplate(userAddress, tokenId, nonce, paymentStr, bidAmount)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary End auction for an NFT - tx template.
// @Description Retrieves tx-template for end auction transaction.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Param tokenId path int true "token id"
// @Param nonce path int true "nonce"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/end-auction/{userAddress}/{tokenId}/{nonce} [get]
func (handler *txTemplateHandler) getEndAuctionTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	nonceStr := c.Param("nonce")

	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.EndAuctionTxTemplate(userAddress, tokenId, nonce)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Deposit EGLD template.
// @Description Retrieves tx-template for deposit transaction.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Param amount path int true "amount"
// @Success 200 {object} formatter.Transaction
// @Router /tx-template/deposit/{userAddress}/{amount} [get]
func (handler *txTemplateHandler) getDepositTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	amountStr := c.Param("amount")

	template := handler.txFormatter.DepositTxTemplate(userAddress, amountStr)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Withdraw EGLD template.
// @Description Retrieves tx-template for withdraw transaction.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Success 200 {object} formatter.Transaction
// @Router /tx-template/withdraw/{userAddress}/{amount} [get]
func (handler *txTemplateHandler) getWithdrawTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	amountStr := c.Param("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.WithdrawTxTemplate(userAddress, amount)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Withdraw Creator Royalties EGLD template.
// @Description Retrieves tx-template for withdraw creator royalties transaction.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Success 200 {object} formatter.Transaction
// @Router /tx-template/withdraw-creator-royalties/{userAddress} [get]
func (handler *txTemplateHandler) getWithdrawCreatorRoyalties(c *gin.Context) {
	userAddress := c.Param("userAddress")

	template := handler.txFormatter.WithdrawCreatorRoyaltiesTxTemplate(userAddress)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Staking Rewards EGLD template.
// @Description Retrieves tx-template for withdraw creator royalties transaction.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path int true "user address"
// @Success 200 {object} formatter.Transaction
// @Router /tx-template/redeem-staking-rewards/{userAddress} [get]
func (handler *txTemplateHandler) getRedeemStakingRewards(c *gin.Context) {
	userAddress := c.Param("userAddress")

	template := handler.txFormatter.RedeemStakingTxTemplate(userAddress)
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
// @Router /tx-template/mint-tokens/{userAddress}/{tokenId}/{numberOfTokens} [get]
func (handler *txTemplateHandler) getMintNftTxTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
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

	if len(collection.ContractAddress) == 0 {
		dtos.JsonResponse(c, http.StatusNotFound, nil, "no contract address")
		return
	}

	//verify collectin is whitelist type, if address is white list, and if have enough to allocated to mint

	//check the type in collection to see if it's a whitelist_minting
	// if collection.Type == entities.Collection_type_whitelisted {

	// 	//get the whitelist by userAddress
	// 	whitelist, errWhitelist := storage.GetWhitelistByAddress(userAddress)

	// 	//check if the whitelist look is invalid
	// 	if errWhitelist != nil {
	// 		if errWhitelist == gorm.ErrRecordNotFound {
	// 			dtos.JsonResponse(c, http.StatusBadRequest, nil, "Address not part of Whitelist")
	// 		} else {
	// 			dtos.JsonResponse(c, http.StatusInternalServerError, nil, errWhitelist.Error())
	// 		}
	// 	}

	// 	//convert numberOfToken string to int
	// 	iNumToken, errConv := strconv.ParseUint(numberOfTokensStr, 10, 64)

	// 	//check conversion error
	// 	if errConv != nil {
	// 		dtos.JsonResponse(c, http.StatusInternalServerError, nil, errConv.Error())
	// 	}

	// 	//have enough to mint,deduct the amount
	// 	if whitelist.Amount >= iNumToken {

	// 		whitelist.Amount -= iNumToken

	// 		errWhitelist = storage.UpdateWhitelist(whitelist)
	// 		if errWhitelist != nil {
	// 			dtos.JsonResponse(c, http.StatusInternalServerError, nil, errWhitelist.Error())

	// 			return
	// 		}
	// 	} else {
	// 		dtos.JsonResponse(c, http.StatusBadRequest, nil, "You are not allocated enough to mint this amount: "+numberOfTokensStr)
	// 	}
	// }

	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	//Signed Message Package

	//get whitelist pem file
	// file, _ := os.Open("config/whitelist-priv.pem")
	// fileBytes, _ := ioutil.ReadAll(file)

	// block, _ := pem.Decode(fileBytes)
	// pkey, erPkey := x509.ParsePKCS8PrivateKey(block.Bytes)

	// if erPkey != nil {
	// 	dtos.JsonResponse(c, http.StatusInternalServerError, nil, erPkey.Error())
	// }

	// edPkey := pkey.(ed25519.PrivateKey)

	// message := tokenId + "" + fmt.Sprint(lastToken.Nonce+1)
	// msg := []byte(message)

	// signedMessage := ed25519.Sign(edPkey, msg)
	if numberOfTokens > 10 {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, "you can't mint more than 10 per batch")
		return
	}
	template, err := handler.txFormatter.NewMintNftsTxTemplate(
		userAddress,
		collection.ContractAddress,
		collection.MintPricePerTokenNominal,
		numberOfTokens,
		collection.CollectionTokenID,
		[]byte(""),
	)
	if err != nil {
		handler.logErr.Println(err.Error())
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, "mint failed")
	}
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Gets tx-template for issue NFT tokens.
// @Description
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Param tokenName path string true "token name"
// @Param tokenTicker path string true "token ticker"
// @Success 200 {object} formatter.Transaction
// @Router /tx-template/issue-nft/{userAddress}/{tokenName}/{tokenTicker} [get]
func (handler *txTemplateHandler) getIssueNFT(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenName := c.Param("tokenName")
	tokenTicker := c.Param("tokenTicker")

	template := handler.txFormatter.NewIssueNFTTxTemplate(
		userAddress,
		tokenName,
		tokenTicker,
	)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Gets tx-template for deploy NFT template contract.
// @Description
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Param tokenId path string true "token id"
// @Param royalties path float64 true "royalties"
// @Param tokenNameBase path string true "tokenNameBase"
// @Param imageBase path string true "imageBase"
// @Param imageExt path string true "imageExt"
// @Param price path float64 true "price"
// @Param maxSupply path int true "maxSupply"
// @Param saleStart path int true "saleStart"
// @Param metadataBase path string true "metadataBase"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/deploy-template/{userAddress}/{tokenId}/{royalties}/{tokenNameBase}/{imageExt}/{price}/{maxSupply}/{saleStart} [get]
func (handler *txTemplateHandler) getDeployNFTTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	tokenNameBase := c.Param("tokenNameBase")
	royaltiesStr := c.Param("royalties")
	imageExt := c.Param("imageExt")
	priceStr := c.Param("price")
	maxSupplyStr := c.Param("maxSupply")
	saleStartStr := c.Param("saleStart")
	imageBase := c.Query("imageBaseLink")
	metadataBase := c.Query("metadataBaseLink")

	maxSupply, err := strconv.ParseUint(maxSupplyStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	saleStart, err := strconv.ParseUint(saleStartStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	royalties, err := strconv.ParseFloat(royaltiesStr, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.DeployNFTTemplateTxTemplate(
		userAddress,
		tokenId,
		royalties,
		tokenNameBase,
		imageBase,
		imageExt,
		price,
		maxSupply,
		saleStart,
		metadataBase,
	)

	//TODO: Grab this response erd SC address
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Gets tx-template for Staking NFT.
// @Description
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Param collectionId path string true "collection id"
// @Param nonce path string true "token nonce"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/deploy-template/{userAddress}/{tokenId}/{royalties}/{tokenNameBase}/{imageExt}/{price}/{maxSupply}/{saleStart} [get]
func (handler *txTemplateHandler) getStakeNftTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	collectionId := c.Param("collectionId")
	nonce := c.Param("nonce")

	intNonce, err := strconv.ParseUint(nonce, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.StakeNFTTemplateTxTemplate(
		userAddress,
		collectionId,
		intNonce,
	)

	//TODO: Grab this response erd SC address
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Gets tx-template for Staking NFT.
// @Description
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Param collectionId path string true "collection id"
// @Param nonce path string true "token nonce"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/deploy-template/{userAddress}/{tokenId}/{royalties}/{tokenNameBase}/{imageExt}/{price}/{maxSupply}/{saleStart} [get]
func (handler *txTemplateHandler) getUnstakeNftTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	collectionId := c.Param("collectionId")
	nonce := c.Param("nonce")

	intNonce, err := strconv.ParseUint(nonce, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template := handler.txFormatter.UnstakeNFTTemplateTxTemplate(
		userAddress,
		collectionId,
		intNonce,
	)

	//TODO: Grab this response erd SC address
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

func (handler *txTemplateHandler) getStakeCollectionTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	collectionId := c.Param("collectionId")

	template := handler.txFormatter.StakeCollectionTemplateTxTemplate(
		userAddress,
		collectionId,
	)

	//TODO: Grab this response erd SC address
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

func (handler *txTemplateHandler) getUnstakeCollectionTemplate(c *gin.Context) {
	userAddress := c.Param("userAddress")
	collectionId := c.Param("collectionId")

	template := handler.txFormatter.UnstakeCollectionTemplateTxTemplate(
		userAddress,
		collectionId,
	)

	//TODO: Grab this response erd SC address
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Gets tx-template for change owner of NFT contract.
// @Description
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Param contractAddress path string true "contract address"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/change-owner/{userAddress}/{contractAddress} [get]
func (handler *txTemplateHandler) getChangeOwner(c *gin.Context) {
	userAddress := c.Param("userAddress")
	contractAddress := c.Param("contractAddress")

	template, err := handler.txFormatter.ChangeOwnerTxTemplate(
		userAddress,
		contractAddress,
	)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Gets tx-template for change set special roles for NFT contract.
// @Description
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Param tokenName path string true "token name"
// @Param tokenTicker path string true "token ticker"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/set-roles/{userAddress}/{tokenId}/{contractAddress} [get]
func (handler *txTemplateHandler) getSetSpecialRoles(c *gin.Context) {
	userAddress := c.Param("userAddress")
	tokenId := c.Param("tokenId")
	contractAddress := c.Param("contractAddress")

	template, err := handler.txFormatter.SetSpecialRolesTxTemplate(
		userAddress,
		tokenId,
		contractAddress,
	)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Gets tx-template for withdraw from Minter SC.
// @Description
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Param contractAddress path string true "contract address"
// @Success 200 {object} formatter.Transaction
// @Router /tx-template/withdraw-minter/{userAddress}/{contractAddress} [get]
func (handler *txTemplateHandler) withdrawFromMinter(c *gin.Context) {
	userAddress := c.Param("userAddress")
	contractAddress := c.Param("contractAddress")

	template := handler.txFormatter.WithdrawFromMinterTxTemplate(userAddress, contractAddress)
	dtos.JsonResponse(c, http.StatusOK, template, "")
}

// @Summary Gets tx-template for request withdraw through Minter.
// @Description The destination will be the Minter Address. Minter will request withdrawal from Marketplace.
// @Tags tx-template
// @Accept json
// @Produce json
// @Param userAddress path string true "user address"
// @Param contractAddress path string true "contract address"
// @Success 200 {object} formatter.Transaction
// @Failure 400 {object} dtos.ApiResponse
// @Router /tx-template/request-withdraw/{userAddress}/{contractAddress} [get]
func (handler *txTemplateHandler) requestWithdrawThroughMinter(c *gin.Context) {
	userAddress := c.Param("userAddress")
	contractAddress := c.Param("contractAddress")

	template, err := handler.txFormatter.RequestWithdrawThroughMinterTxTemplate(userAddress, contractAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, template, "")
}

func (handler *txTemplateHandler) updateSaleStart(c *gin.Context) {
	userAddress := c.Param("userAddress")
	contractAddress := c.Param("contractAddress")
	saleStartStr := c.Param("saleStart")

	/*
		fmt.Println("contract Address: ", contractAddress)

		contractAddressFromBech32, _ := erdgoData.NewAddressFromBech32String(contractAddress)
		contractAddressHex := hex.EncodeToString(contractAddressFromBech32.AddressBytes())

		fmt.Println("newAddress2: " + contractAddressHex)
	*/

	saleStart, err := strconv.ParseUint(saleStartStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template, err := handler.txFormatter.UpdateSaleStartTemplateTxTemplate(userAddress, contractAddress, saleStart)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, template, "")
}

func (handler *txTemplateHandler) updateBuyerWhiteListCheck(c *gin.Context) {
	userAddress := c.Param("userAddress")
	contractAddress := c.Param("contractAddress")
	whiteListCheckStr := c.Param("whiteListCheck")

	/*
		fmt.Println("contract Address: ", contractAddress)

		contractAddressFromBech32, _ := erdgoData.NewAddressFromBech32String(contractAddress)
		contractAddressHex := hex.EncodeToString(contractAddressFromBech32.AddressBytes())

		fmt.Println("newAddress2: " + contractAddressHex)
	*/

	whiteListCheck, err := strconv.ParseUint(whiteListCheckStr, 10, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template, err := handler.txFormatter.UpdateBuyerWhiteListCheckTemplateTxTemplate(userAddress, contractAddress, whiteListCheck)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, template, "")
}

func (handler *txTemplateHandler) getBuyerWhiteListCheck(c *gin.Context) {
	userAddress := c.Param("userAddress")
	contractAddress := c.Param("contractAddress")

	template, err := handler.txFormatter.GetBuyerWhiteListCheckTemplateTxTemplate(userAddress, contractAddress)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, template, "")
}

func (handler *txTemplateHandler) payCheckoutFromDeposit(c *gin.Context) {
	userAddress := c.Param("userAddress")
	amountStr := c.Param("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		dtos.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	template, err := handler.txFormatter.PayCheckoutTxTemplate(userAddress, amount)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	dtos.JsonResponse(c, http.StatusOK, template, "")
}
