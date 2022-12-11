package gui

import (
	"log"
	"net/http"
)

func serveIndexHtml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/index.html")
}

func onAddNewMessage(w http.ResponseWriter, r *http.Request) {
	// TODO:

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func StartGUIServer() *http.Server {
	server := &http.Server{
		Addr: "localhost:5000",
	}

	http.HandleFunc("/", serveIndexHtml)
	http.HandleFunc("/api/messages/new", onAddNewMessage)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return server
}
