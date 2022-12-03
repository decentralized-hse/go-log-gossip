package main

import (
	"github.com/decentralized-hse/go-log-gossip/api"
	"log"
	"os"
	"os/signal"
)

func main() {
	stopApiServerChannel := make(chan bool, 1)
	newLogChannel := make(chan string, 1)

	apiServerConfiguration := api.NewApiServerConfiguration(
		":5001",
		stopApiServerChannel,
		newLogChannel,
	)

	go api.StartApiServer(apiServerConfiguration)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	isNodeAlive := true
	for isNodeAlive {
		select {
		case <-stop:
			stopApiServerChannel <- true
			isNodeAlive = false
		case message := <-newLogChannel:
			log.Println(message)
		}
	}
}
