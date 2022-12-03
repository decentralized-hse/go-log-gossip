package api

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/api/handlers"
	"log"
	"net/http"
)

func ServeAPIServer(configuration *ServerConfiguration) {
	defer configuration.WaitGroup.Done()

	server := &http.Server{
		Addr: configuration.Addr,
	}
	http.HandleFunc("/api/messages/new", handlers.HandleNewLog)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Printf("Started http api server, addr %v", server.Addr)

	<-configuration.Context.Done()
	log.Printf("Shutting down api server")

	_ = server.Shutdown(context.Background())
	log.Printf("Api server sexited properly")
}
