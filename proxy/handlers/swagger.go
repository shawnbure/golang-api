package handlers

import (
	"net/http"

	"github.com/erdsea/erdsea-api/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	baseSwaggerEndpoint = "/swagger"
	anySwaggerEndpoint  = "/*any"
)

func NewSwaggerHandler(groupHandler *groupHandler, conf config.SwaggerConfig) {
	url := ginSwagger.URL(conf.LocalDocRoute)
	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: anySwaggerEndpoint, HandlerFunc: ginSwagger.WrapHandler(swaggerFiles.Handler, url)},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseSwaggerEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	if conf.Enabled {
		groupHandler.AddEndpointGroupHandler(endpointGroupHandler)
	}
}
