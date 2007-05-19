package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/process"
	"github.com/erdsea/erdsea-api/proxy/handlers"
	"github.com/erdsea/erdsea-api/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/erdsea/erdsea-api/docs"
)

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

	url := ginSwagger.URL("http://localhost:5000/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

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
