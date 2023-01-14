package commands

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/storage"
)

type GetLogsHandler struct {
	storage storage.LogStorage
}

func NewGetLogsHandler(storage storage.LogStorage) *GetLogsHandler {
	return &GetLogsHandler{storage: storage}
}

func (c *GetLogsHandler) Handle(_ context.Context, command *GetLogsQuery) (response *GetLogsResponse, err error) {
	logs, err := c.storage.GetNodeLogs(command.NodeId)
	if err != nil {
		return nil, err
	}

	response = &GetLogsResponse{Logs: logs}
	return
}
