package start

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/api"
	"github.com/decentralized-hse/go-log-gossip/domain/features/logs/commands"
	"github.com/decentralized-hse/go-log-gossip/infra/config"
	"github.com/decentralized-hse/go-log-gossip/infra/gossip"
	"github.com/decentralized-hse/go-log-gossip/infra/keys"
	"github.com/decentralized-hse/go-log-gossip/storage"
	"github.com/mehdihadeli/go-mediatr"
	"log"
	"os"
	"os/signal"
	"sync"
)

var (
	wg          sync.WaitGroup
	ctx, cancel = context.WithCancel(context.Background())
	cfg         *config.Config
	keysPair    *keys.PublicPrivateKeyPair
	gossiper    *gossip.Gossiper
	logStorage  *storage.InMemoryStorage
	selfNodeId  string
)

func CommandStartNode() error {
	loadConfig()
	loadKeys()
	startGossip()
	initializeLogStorage()
	initializeMediatr()
	startApiServer()
	registerGracefulShutdown()

	return nil
}

func loadConfig() {
	loadedCfg, err := config.LoadFromFile()
	if err != nil {
		panic(err)
	}
	cfg = loadedCfg
}

func loadKeys() {
	loadFromFiles, err := keys.LoadFromFiles(cfg.Paths.PathToFolderWithRSAKeys)
	if err != nil {
		panic(err)
	}
	keysPair = loadFromFiles

	selfNodeId = keysPair.GetPublicKey().Encode()
}

func initializeMediatr() {
	createSelfLogHandler := commands.NewCreateSelfLogHandler(logStorage, gossiper, selfNodeId)
	addLogHandler := commands.NewAddLogHandler(logStorage, gossiper)
	getLogsHandler := commands.NewGetLogsHandler(logStorage)
	sendLogHandler := commands.NewSendLogHandler(logStorage, gossiper)

	err := mediatr.RegisterRequestHandler[*commands.CreateSelfLogCommand, *commands.CreateSelfLogResponse](createSelfLogHandler)
	if err != nil {
		panic("Failed to register CreateSelfLogHandler")
	}

	err = mediatr.RegisterRequestHandler[*commands.AddLogCommand, *commands.AddLogResponse](addLogHandler)
	if err != nil {
		panic("Failed to register AddLogHandler")
	}

	err = mediatr.RegisterRequestHandler[*commands.GetLogsQuery, *commands.GetLogsResponse](getLogsHandler)
	if err != nil {
		panic("Failed to register GetLogsHandler")
	}

	err = mediatr.RegisterRequestHandler[*commands.SendLogCommand, *commands.SendLogResponse](sendLogHandler)
	if err != nil {
		panic("Failed to register GetLogsHandler")
	}
}

func startApiServer() {
	apiServerConfiguration := api.NewServerConfiguration(
		cfg.Api.Addr,
		ctx,
		&wg,
	)

	go api.ServeAPIServer(apiServerConfiguration)
	wg.Add(1)
}

func startGossip() {
	gossiper = gossip.Start(cfg, keysPair, gossipHandler, ctx, &wg)
	wg.Add(1)
}

func initializeLogStorage() {
	logStorage = storage.NewInMemoryStorage()
}

func registerGracefulShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	log.Println("Graceful shutdown. Cancelling all goroutines")
	cancel()

	log.Println("Send cancel signal, waiting goroutines")
	wg.Wait()
}
