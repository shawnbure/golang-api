package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/erdsea/erdsea-api/config"
	_ "github.com/erdsea/erdsea-api/docs"
	"github.com/erdsea/erdsea-api/process"
	"github.com/erdsea/erdsea-api/proxy/handlers"
	"github.com/erdsea/erdsea-api/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var corsHeaders = []string{
	"Origin",
	"Content-Length",
	"Content-Type",
	"Authorization",
}

type webServer struct {
	router        *gin.Engine
	generalConfig *config.GeneralConfig
}

// @title erdsea-api
// @version 1.0
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
func NewWebServer(cfg *config.GeneralConfig) (*webServer, error) {
	router := gin.Default()
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowHeaders = corsHeaders
	corsCfg.AllowAllOrigins = true
	router.Use(cors.New(corsCfg))

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
	handlers.NewTokensHandler(groupHandler)
	handlers.NewCollectionsHandler(groupHandler, cfg.Auth, cfg.Blockchain)
	handlers.NewTransactionsHandler(groupHandler)
	handlers.NewTxTemplateHandler(groupHandler, cfg.Auth, cfg.Blockchain)
	handlers.NewPriceHandler(groupHandler, cfg.Auth)
	handlers.NewAccountsHandler(groupHandler, cfg.Auth)
	handlers.NewSearchHandler(groupHandler, cfg.Auth)
	handlers.NewSwaggerHandler(groupHandler, cfg.Swagger)

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
