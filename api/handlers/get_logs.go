package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain/features/logs/commands"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type getLogRequest struct {
	NodeId string `json:"node_id"`
}

func HandleGetLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var getLogMessage getLogRequest
	err := json.NewDecoder(r.Body).Decode(&getLogMessage)

	if err != nil {
		http.Error(w, "Error while reading body", http.StatusBadRequest)
		return
	}

	response, err := mediatr.Send[*commands.GetLogsQuery, *commands.GetLogsResponse](
		r.Context(),
		commands.NewGetLogsCommand(getLogMessage.NodeId),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while trying to get logs by id=%s", getLogMessage.NodeId), http.StatusInternalServerError)
		return
	}

	logsJson, err := json.Marshal(response.Logs)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(logsJson)
	if err != nil {
		panic(err)
	}
}
