package commands

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/domain/features/dtos"
	"github.com/decentralized-hse/go-log-gossip/storage"
)

type GetLogsHandler struct {
	storage storage.LogStorage
}

func NewGetLogsHandler(storage storage.LogStorage) *GetLogsHandler {
	return &GetLogsHandler{storage: storage}
}

func (c *GetLogsHandler) Handle(_ context.Context, command *GetLogsQuery) (*GetLogsResponse, error) {
	logs, err := c.storage.GetNodeLogs(command.NodeId)
	if err != nil {
		return nil, err
	}

	dtoLogs := make([]*dtos.LogDTO, len(logs))
	for i, log := range logs {
		dtoLogs[i] = dtos.NewLogDTO(log)
	}

	return &GetLogsResponse{Logs: dtoLogs}, nil
}
