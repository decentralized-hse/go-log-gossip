package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/domain"
	"github.com/decentralized-hse/go-log-gossip/domain/features/creating_self_log/commands"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
	"strings"
)

func HandleGetLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathParams := strings.Split(r.URL.Path, "/")
	nodeId := pathParams[len(pathParams)-1]

	response, err := mediatr.Send[*commands.GetLogsCommand, *commands.GetLogsResponse](
		r.Context(),
		commands.NewGetLogsCommand(domain.NodeId(nodeId)),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while trying to get logs by id=%s", nodeId), http.StatusInternalServerError)
		return
	}

	logsJson, err := json.Marshal(response.Logs)

	_, err = w.Write(logsJson)
	if err != nil {
		panic(err)
	}
}
