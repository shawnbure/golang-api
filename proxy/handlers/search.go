package handlers

import (
	"net/http"

	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/services"
	"github.com/gin-gonic/gin"
)

const (
	baseSearchEndpoint        = "/search"
	generalSearchEndpoint     = "/:searchString"
	collectionsSearchEndpoint = "/collections/:collectionName"
	accountsSearchEndpoint    = "/accounts/:accountName"

	SearchCategoryLimit = 5
)

type GeneralSearchResponse struct {
	Accounts    []entities.Account
	Collections []entities.Collection
}

type searchHandler struct {
}

func NewSearchHandler(groupHandler *groupHandler) {
	handler := &searchHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: generalSearchEndpoint, HandlerFunc: handler.search},
		{Method: http.MethodGet, Path: collectionsSearchEndpoint, HandlerFunc: handler.collectionSearch},
		{Method: http.MethodGet, Path: accountsSearchEndpoint, HandlerFunc: handler.accountSearch},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseSearchEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

// @Summary General search by string.
// @Description Searches for collections by name and accounts by name. Cached for 20 minutes. Limit 5 elements for each.
// @Tags search
// @Accept json
// @Produce json
// @Param searchString path string true "search string"
// @Success 200 {object} GeneralSearchResponse
// @Failure 500 {object} dtos.ApiResponse
// @Router /search/{searchString} [get]
func (handler *searchHandler) search(c *gin.Context) {
	searchString := c.Param("searchString")

	collections, err := services.GetCollectionsWithNameAlike(searchString, SearchCategoryLimit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	accounts, err := services.GetAccountsWithNameAlike(searchString, SearchCategoryLimit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response := GeneralSearchResponse{
		Accounts:    accounts,
		Collections: collections,
	}
	dtos.JsonResponse(c, http.StatusOK, response, "")
}

// @Summary Search collections by name.
// @Description Searches for collections by name. Cached for 20 minutes. Limit 5 elements.
// @Tags search
// @Accept json
// @Produce json
// @Param collectionName path string true "search string"
// @Success 200 {object} []entities.Collection
// @Failure 500 {object} dtos.ApiResponse
// @Router /search/collections/{collectionName} [get]
func (handler *searchHandler) collectionSearch(c *gin.Context) {
	collectionName := c.Param("collectionName")

	collections, err := services.GetCollectionsWithNameAlike(collectionName, SearchCategoryLimit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, collections, "")
}

// @Summary Search accounts by name.
// @Description Searches for accounts by name. Cached for 20 minutes. Limit 5 elements.
// @Tags search
// @Accept json
// @Produce json
// @Param accountName path string true "search string"
// @Success 200 {object} []entities.Account
// @Failure 500 {object} dtos.ApiResponse
// @Router /search/accounts/{accountName} [get]
func (handler *searchHandler) accountSearch(c *gin.Context) {
	accountName := c.Param("accountName")

	accounts, err := services.GetAccountsWithNameAlike(accountName, SearchCategoryLimit)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	dtos.JsonResponse(c, http.StatusOK, accounts, "")
}
