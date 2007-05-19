package proxy

import (
	"context"
	"net/http"
	"strings"

	"github.com/erdsea/erdsea-api/alerts/tg"
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

var ctx = context.Background()

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

	bot, err := makeBot(cfg.Bot)
	if err != nil {
		return nil, err
	}

	observerMonitor := process.NewObserverMonitor(
		bot,
		ctx,
		cfg.Monitor.ObserverMonitorEnable,
	)

	processor := process.NewEventProcessor(
		cfg.ConnectorApi.Addresses,
		cfg.ConnectorApi.Identifiers,
		observerMonitor,
	)

	err = handlers.NewEventsHandler(
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

	handlers.NewAuthHandler(groupHandler, *authService)
	handlers.NewTokensHandler(groupHandler)
	handlers.NewCollectionsHandler(groupHandler, cfg.Auth, cfg.Blockchain)
	handlers.NewTransactionsHandler(groupHandler)
	handlers.NewTxTemplateHandler(groupHandler, cfg.Blockchain)
	handlers.NewPriceHandler(groupHandler)
	handlers.NewAccountsHandler(groupHandler, cfg.Auth)
	handlers.NewSearchHandler(groupHandler)
	handlers.NewSwaggerHandler(groupHandler, cfg.Swagger)

	groupHandler.RegisterEndpoints(router)

	return &webServer{
		router:        router,
		generalConfig: cfg,
	}, nil
}

func (w *webServer) Run() *http.Server {
	address := w.generalConfig.ConnectorApi.Address
	if !strings.Contains(address, ":") {
		panic("bad address")
	}
	server := &http.Server{
		Addr:    address,
		Handler: w.router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return server
}

func makeBot(cfg config.BotConfig) (tg.Bot, error) {
	if !cfg.Enable {
		return &tg.DisabledBot{}, nil
	}

	return tg.NewTelegramBot(cfg)
}
