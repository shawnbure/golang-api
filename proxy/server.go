package proxy

import (
	"fmt"
	"github.com/erdsea/erdsea-api/process"
	"net/http"
	"strings"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/proxy/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type webServer struct {
	router        *gin.Engine
	generalConfig *config.GeneralConfig
}

func NewWebServer(generalConfig *config.GeneralConfig) (*webServer, error) {
	router := gin.Default()
	router.Use(cors.Default())

	groupHandler := handlers.NewGroupHandler()

	processor := process.NewEventProcessor(nil, nil)
	err := handlers.NewEventsHandler(groupHandler, processor, generalConfig.ConnectorApi)
	if err != nil {
		return nil, err
	}

	handlers.NewAssetsHandler(groupHandler)

	groupHandler.RegisterEndpoints(router)

	return &webServer{
		router:        router,
		generalConfig: generalConfig,
	}, nil
}

func (w *webServer) Run() *http.Server {
	port := w.generalConfig.ConnectorApi.Port
	if !strings.Contains(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}
	server := &http.Server{
		Addr:    port,
		Handler: w.router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return server
}
