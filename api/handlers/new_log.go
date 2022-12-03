package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func OnNewLogApiHandlerFactory(newRawLogChannel chan<- string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handleNewLogHandler(newRawLogChannel, w, r)
	}
}

type newLogMessage struct {
	Log string `json:"log"`
}

func handleNewLogHandler(newRawLogChannel chan<- string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logMessage newLogMessage
	err := json.NewDecoder(r.Body).Decode(&logMessage)

	if err != nil {
		http.Error(w, "Error while reading body", http.StatusBadRequest)
		return
	}
	log.Println(logMessage)

	newRawLogChannel <- logMessage.Log

	_, err = w.Write([]byte("OK"))
	if err != nil {
		panic(err)
	}
}
