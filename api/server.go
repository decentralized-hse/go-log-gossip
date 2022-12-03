package api

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/api/handlers"
	"log"
	"net/http"
	"time"
)

func StartApiServer(configuration *ApiServerConfiguration) {
	newMessageHandler := handlers.OnNewLogApiHandlerFactory(configuration.newLogChannel)
	http.HandleFunc("/api/messages/new", newMessageHandler)

	server := &http.Server{
		Addr: configuration.addr,
	}

	log.Printf(server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Printf("Started server")

	<-configuration.channelToOffServer

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %+v", err)
	}
	log.Printf("Server exited properly")
}
