package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/cdn"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/interaction"
	"github.com/erdsea/erdsea-api/logging"
	"github.com/erdsea/erdsea-api/proxy"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/urfave/cli"
)

const (
	defaultLogsPath    = "logs"
	logFilePrefix      = "erdsea"
	logFileLifeSpanSec = 86400
)

var (
	backgroundContextTimeout = 5 * time.Second
)

var (
	cliHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`
	log = logger.GetOrCreate("erdsea-api")

	logLevel = cli.StringFlag{
		Name:  "log-level",
		Usage: "This flag specifies the log level. Options: *:NONE | ERROR | WARN | INFO | DEBUG | TRACE",
		Value: fmt.Sprintf("*:%s", logger.LogInfo.String()),
	}

	logSaveFile = cli.BoolFlag{
		Name:  "log-save",
		Usage: "Boolean option for enabling log saving",
	}

	generalConfigFile = cli.StringFlag{
		Name:  "general-config",
		Usage: "The path for the general config",
		Value: "./config/config.toml",
	}

	workingDirectory = cli.StringFlag{
		Name:  "working-directory",
		Usage: "This flag specifies the directory where the proxy will store logs.",
		Value: "",
	}
)

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = cliHelpTemplate
	app.Name = "erdsea-api"
	app.Flags = []cli.Flag{
		logLevel,
		logSaveFile,
		generalConfigFile,
		workingDirectory,
	}
	app.Authors = []cli.Author{
		{
			Name:  "The ErdSea Team",
			Email: "contact@erdsea.com",
		},
	}
	app.Action = startProxy

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
}

func startProxy(ctx *cli.Context) error {
	log.Info("starting erdsea-api proxy...")

	fileLogging, err := initLogger(ctx)
	if err != nil {
		return err
	}

	generalConfigPath := ctx.GlobalString(generalConfigFile.Name)
	cfg, err := config.LoadConfig(generalConfigPath)
	if err != nil {
		return err
	}

	establishConnections(cfg)

	api, err := proxy.NewWebServer(cfg)
	if err != nil {
		return err
	}

	server := api.Run()

	waitForGracefulShutdown(server)
	log.Debug("closing erdsea-api proxy...")
	if !check.IfNil(fileLogging) {
		err = fileLogging.Close()
		if err != nil {
			return err
		}
	}

	cache.CloseCacher()

	return nil
}

func establishConnections(cfg *config.GeneralConfig) {
	interaction.InitBlockchainInteractor(cfg.Blockchain)
	cache.InitCacher(cfg.Cache)
	storage.Connect(cfg.Database)
	cdn.MakeCloudyCDN(cfg.CDN)
}

func initLogger(ctx *cli.Context) (logging.FileLogger, error) {
	logLevelValue := ctx.GlobalString(logLevel.Name)

	err := logger.SetLogLevel(logLevelValue)
	if err != nil {
		return nil, err
	}

	workingDir, err := getWorkingDir(ctx)
	if err != nil {
		return nil, err
	}

	var fileLogging logging.FileLogger
	saveLogs := ctx.GlobalBool(logSaveFile.Name)
	if saveLogs {
		fileLogging, err = logging.NewFileLogging(workingDir, defaultLogsPath, logFilePrefix)
		if err != nil {
			return fileLogging, err
		}
	}

	if !check.IfNil(fileLogging) {
		err = fileLogging.ChangeFileLifeSpan(time.Second * time.Duration(logFileLifeSpanSec))
		if err != nil {
			return nil, err
		}
	}

	return fileLogging, nil
}

func waitForGracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), backgroundContextTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
	_ = server.Close()
}

func getWorkingDir(ctx *cli.Context) (string, error) {
	if ctx.IsSet(workingDirectory.Name) {
		return ctx.GlobalString(workingDirectory.Name), nil
	}

	return os.Getwd()
}
