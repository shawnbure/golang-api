package handlers

import (
	"fmt"
	"net/http"

	"github.com/erdsea/erdsea-api/cdn"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/gin-gonic/gin"
)

const (
	baseImageEndpoint = "/image"
	getImageEndpoint  = "/:filename"

	contentTypeImage = "image/%s"
)

type imageHandler struct{}

func NewImageHandler(groupHandler *groupHandler) {
	handler := &imageHandler{}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: getImageEndpoint, HandlerFunc: handler.getImage},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseImageEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
}

func (h *imageHandler) getImage(c *gin.Context) {
	filename := c.Param("filename")

	uploader, err := cdn.GetImageUploaderOrErr()
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	img, imgType, err := uploader.GetImage(filename)
	if err != nil {
		dtos.JsonResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	c.Data(http.StatusOK, fmt.Sprintf(contentTypeImage, imgType), img)
}
