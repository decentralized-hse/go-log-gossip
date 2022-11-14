package main

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/gui"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Printf("main: starting HTTP server")

	server := gui.StartGUIServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
	}
}
