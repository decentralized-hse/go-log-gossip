package handlers

import (
	"encoding/json"
	"github.com/decentralized-hse/go-log-gossip/domain/features/creating_self_log/commands"
	"github.com/mehdihadeli/go-mediatr"
	"log"
	"net/http"
)

type newLogMessage struct {
	Log string `json:"log"`
}

func HandleNewLog(w http.ResponseWriter, r *http.Request) {
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

	response, err := mediatr.Send[*commands.CreateSelfLogCommand, *commands.CreateSelfLogResponse](
		r.Context(),
		commands.NewCreateSelfLogCommand(logMessage.Log),
	)
	if err != nil {
		http.Error(w, "Error while trying to add new log", http.StatusInternalServerError)
		return
	}

	log.Println(response)

	_, err = w.Write([]byte("OK"))
	if err != nil {
		panic(err)
	}
}
