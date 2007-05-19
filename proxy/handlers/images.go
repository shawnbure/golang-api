package handlers

import (
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/proxy/middleware"
	"github.com/erdsea/erdsea-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	baseImagesEndpoint           = "/images"

	getAccountProfileEndpoint    = "/get/account/profile/:userAddress"
	getAccountCoverEndpoint      = "/get/account/cover/:userAddress"

	setAccountProfileEndpoint    = "/set/account/profile/:userAddress"
	setAccountCoverEndpoint      = "/set/account/cover/:userAddress"

	getCollectionProfileEndpoint = "/get/collection/profile/:collectionName"
	getCollectionCoverEndpoint   = "/get/collection/cover/:collectionName"

	setCollectionProfileEndpoint = "/set/collection/profile/:collectionName"
	setCollectionCoverEndpoint   = "/set/collection/cover/:collectionName"
)

type imageHandler struct {
}

func NewImageHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &imageHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: getAccountProfileEndpoint, HandlerFunc: handler.getAccountProfile},
		{Method: http.MethodPost, Path: setAccountProfileEndpoint, HandlerFunc: handler.setAccountProfile},


		{Method: http.MethodGet, Path: getAccountCoverEndpoint, HandlerFunc: handler.getAccountCover},
		{Method: http.MethodPost, Path: setAccountCoverEndpoint, HandlerFunc: handler.setAccountCover},


		{Method: http.MethodGet, Path: getCollectionProfileEndpoint, HandlerFunc: handler.getCollectionProfile},
		{Method: http.MethodPost, Path: setCollectionProfileEndpoint, HandlerFunc: handler.setCollectionProfile},


		{Method: http.MethodGet, Path: getCollectionCoverEndpoint, HandlerFunc: handler.getCollectionCover},
		{Method: http.MethodPost, Path: setCollectionCoverEndpoint, HandlerFunc: handler.setCollectionCover},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseImagesEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *imageHandler) getAccountProfile(c *gin.Context) {
	userAddress := c.Param("userAddress")

	image, err := services.GetAccountProfileImage(userAddress)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, &image, "")
}

func (handler *imageHandler) setAccountProfile(c *gin.Context) {
	var imageBase64 string
	userAddress := c.Param("userAddress")

	err := c.Bind(&imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress, exists := c.Get(middleware.AddressKey)
	if !exists || jwtAddress != userAddress {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	err = services.SetAccountProfileImage(userAddress, &imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, "", "")
}

func (handler *imageHandler) getAccountCover(c *gin.Context) {
	userAddress := c.Param("userAddress")

	image, err := services.GetAccountCoverImage(userAddress)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, &image, "")
}

func (handler *imageHandler) setAccountCover(c *gin.Context) {
	var imageBase64 string
	userAddress := c.Param("userAddress")

	err := c.Bind(&imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	jwtAddress, exists := c.Get(middleware.AddressKey)
	if !exists || jwtAddress != userAddress {
		data.JsonResponse(c, http.StatusUnauthorized, nil, "")
		return
	}

	err = services.SetAccountCoverImage(userAddress, &imageBase64)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, "", "")
}

func (handler *imageHandler) getCollectionProfile(c *gin.Context) {

}

func (handler *imageHandler) setCollectionProfile(c *gin.Context) {

}

func (handler *imageHandler) getCollectionCover(c *gin.Context) {

}

func (handler *imageHandler) setCollectionCover(c *gin.Context) {

}
