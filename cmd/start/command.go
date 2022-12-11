package start

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/api"
	"github.com/decentralized-hse/go-log-gossip/domain/features/creating_self_log/commands"
	"github.com/decentralized-hse/go-log-gossip/infra/config"
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
)

func CommandStartNode() error {
	loadConfig()
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

func initializeMediatr() {
	createSelfLogHandler := commands.NewCreateSelfLogHandler(nil)
	err := mediatr.RegisterRequestHandler[*commands.CreateSelfLogCommand, *commands.CreateSelfLogResponse](createSelfLogHandler)
	if err != nil {
		panic("Failed to register CreateSelfLogHandler")
	}
}

func startApiServer() {
	apiServerConfiguration := api.NewServerConfiguration(
		":5001",
		ctx,
		&wg,
	)

	wg.Add(1)
	go api.ServeAPIServer(apiServerConfiguration)
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
