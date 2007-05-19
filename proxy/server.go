package proxy

import (
	"fmt"
	"github.com/erdsea/erdsea-api/services"
	"net/http"
	"strings"

	"github.com/erdsea/erdsea-api/process"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/proxy/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type webServer struct {
	router        *gin.Engine
	generalConfig *config.GeneralConfig
}

func NewWebServer(cfg *config.GeneralConfig) (*webServer, error) {
	router := gin.Default()
	router.Use(cors.Default())

	groupHandler := handlers.NewGroupHandler()

	processor := process.NewEventProcessor(
		cfg.ConnectorApi.Addresses,
		cfg.ConnectorApi.Identifiers,
	)

	err := handlers.NewEventsHandler(
		groupHandler,
		processor,
		cfg.ConnectorApi,
	)
	if err != nil {
		return nil, err
	}

	authService, err := services.NewAuthService(cfg.Auth)
	if err != nil {
		return nil, err
	}

	//TODO: think about handlers params - maybe single param
	handlers.NewAuthHandler(groupHandler, *authService)
	handlers.NewAssetsHandler(groupHandler, cfg.Auth)
	handlers.NewCollectionsHandler(groupHandler, cfg.Auth, cfg.Blockchain)
	handlers.NewTransactionsHandler(groupHandler, cfg.Auth)
	handlers.NewTxTemplateHandler(groupHandler, cfg.Auth, cfg.Blockchain)
	handlers.NewPriceHandler(groupHandler, cfg.Auth)

	groupHandler.RegisterEndpoints(router)

	return &webServer{
		router:        router,
		generalConfig: cfg,
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
