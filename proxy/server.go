package proxy

import (
	"context"
	"net/http"
	"strings"

	"github.com/ENFT-DAO/youbei-api/alerts/tg"
	"github.com/ENFT-DAO/youbei-api/config"
	_ "github.com/ENFT-DAO/youbei-api/docs"
	"github.com/ENFT-DAO/youbei-api/indexer"
	"github.com/ENFT-DAO/youbei-api/process"
	"github.com/ENFT-DAO/youbei-api/proxier"
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
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowHeaders = corsHeaders
	corsCfg.AllowAllOrigins = true
	//Access-Control-Allow-Origin
	//corsCfg.a

	router.Use(cors.New(corsCfg))

	groupHandler := handlers.NewGroupHandler()

	bot, err := makeBot(cfg.Bot)
	if err != nil {
		return nil, err
	}
	marketPlaceIndexer, err := indexer.NewMarketPlaceIndexer(cfg.Blockchain.MarketplaceAddress, cfg.Blockchain.ApiUrl, cfg.Blockchain.ApiUrlSec, cfg.Blockchain.CollectionAPIDelay)
	if err != nil {
		return nil, err
	}
	proxier.SetIPs(cfg.Proxy.List)
	collectionIndexer, err := indexer.NewCollectionIndexer(cfg.Blockchain.DeployerAddress, cfg.Blockchain.ApiUrl, cfg.Blockchain.ApiUrlSec, cfg.Blockchain.CollectionAPIDelay)
	if err != nil {
		return nil, err
	}
	go collectionIndexer.StartWorker()
	go marketPlaceIndexer.StartWorker()
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
	handlers.NewTokensHandler(groupHandler, cfg.Auth, cfg.Blockchain)
	handlers.NewCollectionsHandler(groupHandler, cfg.Auth, cfg.Blockchain, cfg.CarbonSetting)

	handlers.NewSessionStatesHandler(groupHandler, cfg.Auth, cfg.Blockchain)

	handlers.NewTransactionsHandler(groupHandler)
	handlers.NewTxTemplateHandler(groupHandler, cfg.Blockchain)
	handlers.NewPriceHandler(groupHandler)
	handlers.NewAccountsHandler(groupHandler, cfg.Auth)
	handlers.NewSearchHandler(groupHandler)
	handlers.NewSwaggerHandler(groupHandler, cfg.Swagger)
	handlers.NewDepositsHandler(groupHandler, cfg.Blockchain)
	handlers.NewRoyaltiesHandler(groupHandler, cfg.Blockchain)
	handlers.NewImageHandler(groupHandler)
	handlers.NewStatsHandler(groupHandler)
	handlers.NewReportHandler(groupHandler)
	handlers.NewActivitiesHandler(groupHandler)
	handlers.NewExplorerHandler(groupHandler)

	handlers.NewDreamshipHandler(groupHandler, cfg.ExternalCredential)

	//

	groupHandler.RegisterEndpoints(router)

	return &webServer{
		router:        router,
		generalConfig: cfg,
	}, nil
}

func (w *webServer) Run() *http.Server {

	// ep := blockchain.NewElrondProxy("https://devnet-gateway.elrond.com", nil)

	// vmRequest := &data.VmValueRequest{
	// 	Address:    "erd1qqqqqqqqqqqqqpgqgmv7tw2a9lvcyxw2wrmskevvc0955cagy4wsrapmzg",
	// 	FuncName:   "getPrice",
	// 	CallerAddr: "erd1p39zv9xw5ftpfxy9s9afzkjaafadk9na44fput904luqgmpmh8rsrtwufq",
	// 	CallValue:  "",
	// 	Args:       nil,
	// }
	// response, err := ep.ExecuteVMQuery(vmRequest)
	// if err != nil {
	// 	log.Fatal("error executing vm query", "error", err)
	// }
	// bytes := response.Data.ReturnData[0]
	// num := big.NewInt(0).SetBytes(bytes)
	// log.Print("response", "contract version", num)

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
