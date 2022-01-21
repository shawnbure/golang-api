package proxy

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/ENFT-DAO/youbei-api/alerts/tg"
	"github.com/ENFT-DAO/youbei-api/config"
	_ "github.com/ENFT-DAO/youbei-api/docs"
	"github.com/ENFT-DAO/youbei-api/process"
	"github.com/ENFT-DAO/youbei-api/proxy/handlers"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var corsHeaders = []string{
	"Origin",
	"Content-Length",
	"Content-Type",
	"Authorization",
	"X-Requested-With",
	"Accept",
	"Accept-Encoding",
	"X-CSRF-Token",
	"Cache-Control",
}

var corsMethods = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"DELETE",
	"OPTIONS",
	"PATCH",
}

var corsOrigins = []string{
	"https://*.youbei.io",
}

var ctx = context.Background()

type webServer struct {
	router        *gin.Engine
	generalConfig *config.GeneralConfig
}

// @title youbei-api
// @version 1.0
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
func NewWebServer(cfg *config.GeneralConfig) (*webServer, error) {

	router := gin.Default()

	//corsCfg := cors.DefaultConfig()
	//corsCfg.AllowWildcard = true //allows widcards in the origin domain
	//corsCfg.AllowHeaders = corsHeaders
	//corsCfg.AllowMethods = corsMethods
	//corsCfg.AllowAllOrigins = true
	//corsCfg.AllowOrigins = corsOrigins

	//router.Use(cors.New(corsCfg))

	router.Use(cors.New(cors.Config{
		AllowWildcard:    true, //allows widcards in the origin domain
		AllowOrigins:     corsOrigins,
		AllowMethods:     corsMethods,
		AllowHeaders:     corsHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		/*
			AllowOriginFunc: func(origin string) bool {
				return origin == "https://github.com"
			},
		*/
		MaxAge: 12 * time.Hour,
	}))

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
		cfg.Blockchain.ProxyUrl,
		cfg.Blockchain.MarketplaceAddress,
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
	handlers.NewTokensHandler(groupHandler, cfg.Blockchain)
	handlers.NewCollectionsHandler(groupHandler, cfg.Auth, cfg.Blockchain)
	handlers.NewTransactionsHandler(groupHandler)
	handlers.NewTxTemplateHandler(groupHandler, cfg.Blockchain)
	handlers.NewPriceHandler(groupHandler)
	handlers.NewAccountsHandler(groupHandler, cfg.Auth)
	handlers.NewSearchHandler(groupHandler)
	handlers.NewSwaggerHandler(groupHandler, cfg.Swagger)
	handlers.NewDepositsHandler(groupHandler, cfg.Blockchain)
	handlers.NewRoyaltiesHandler(groupHandler, cfg.Blockchain)
	handlers.NewImageHandler(groupHandler)

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
