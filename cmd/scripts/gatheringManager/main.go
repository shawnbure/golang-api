package main

import (
	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/cdn"
	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/interaction"
	"github.com/ENFT-DAO/youbei-api/stats/gatherer"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/urfave/cli"
	"os"
	"os/signal"
)

var (
	generalConfigFile = cli.StringFlag{
		Name:  "general-config",
		Usage: "The path for the general config",
		Value: getWorkingDirectory("config/config.toml"),
	}

	workingDirectory = cli.StringFlag{
		Name:  "working-directory",
		Usage: "This flag specifies the directory where the proxy will store logs.",
		Value: "",
	}
)

func main() {
	app := cli.NewApp()

	app.Name = "youbei-api"
	app.Action = startProxy
	app.Flags = []cli.Flag{
		generalConfigFile,
		workingDirectory,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func getWorkingDirectory(param string) string {
	dir, dir_err := os.Getwd()
	if dir_err != nil {
		panic(dir_err)
	}
	return dir + "/" + param
}

func startProxy(ctx *cli.Context) error {
	generalConfigPath := ctx.GlobalString(generalConfigFile.Name)
	cfg, err := config.LoadConfig(generalConfigPath)
	if err != nil {
		return err
	}

	establishConnections(cfg)

	agg := gatherer.GetManager()
	if agg != nil {
		agg.Start()
	}

	waitForGracefulShutdown()
	cache.CloseCacher()
	return nil
}

func establishConnections(cfg *config.GeneralConfig) {
	interaction.InitBlockchainInteractor(cfg.Blockchain)
	cache.InitCacher(cfg.Cache)
	storage.Connect(cfg.Database)
	cdn.InitUploader(cfg.CDN)
}

func waitForGracefulShutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	agg := gatherer.GetManager()
	if agg != nil {
		agg.Stop()
	}
}

func getWorkingDir(ctx *cli.Context) (string, error) {
	if ctx.IsSet(workingDirectory.Name) {
		return ctx.GlobalString(workingDirectory.Name), nil
	}

	return os.Getwd()
}
