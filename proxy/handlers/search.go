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
	baseSearchEndpoint        = "/search"
	generalSearchEndpoint     = "/:searchString"
	collectionsSearchEndpoint = "/collections/:collectionName"
	accountsSearchEndpoint    = "/accounts/:accountName"

	SearchCategoryLimit = 5
)

type generalSearchResponse struct {
	accounts    []data.Account
	collections []data.Collection
}

type searchHandler struct {
}

func NewSearchHandler(groupHandler *groupHandler, authCfg config.AuthConfig) {
	handler := &searchHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: generalSearchEndpoint, HandlerFunc: handler.search},
		{Method: http.MethodGet, Path: collectionsSearchEndpoint, HandlerFunc: handler.collectionSearch},
		{Method: http.MethodGet, Path: accountsSearchEndpoint, HandlerFunc: handler.accountSearch},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseSearchEndpoint,
		Middlewares:      []gin.HandlerFunc{middleware.Authorization(authCfg.JwtSecret)},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (handler *searchHandler) search(c *gin.Context) {
	searchString := c.Param("searchString")

	collections, err := services.GetCollectionsWithNameAlike(searchString, SearchCategoryLimit)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	accounts, err := services.GetAccountsWithNameAlike(searchString, SearchCategoryLimit)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response := generalSearchResponse{
		accounts:    accounts,
		collections: collections,
	}
	data.JsonResponse(c, http.StatusOK, response, "")
}

func (handler *searchHandler) collectionSearch(c *gin.Context) {
	collectionName := c.Param("collectionName")

	collections, err := services.GetCollectionsWithNameAlike(collectionName, SearchCategoryLimit)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, collections, "")
}

func (handler *searchHandler) accountSearch(c *gin.Context) {
	accountName := c.Param("accountName")

	accounts, err := services.GetAccountsWithNameAlike(accountName, SearchCategoryLimit)
	if err != nil {
		data.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	data.JsonResponse(c, http.StatusOK, accounts, "")
}
