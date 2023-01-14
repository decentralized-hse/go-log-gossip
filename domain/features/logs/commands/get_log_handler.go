package commands

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/infra/gossip"
	"github.com/decentralized-hse/go-log-gossip/storage"
)

type GetLogHandler struct {
	storage  storage.LogStorage
	gossiper gossip.Gossiper
}

func NewGetLogHandler(storage storage.LogStorage, gossiper gossip.Gossiper) *GetLogHandler {
	return &GetLogHandler{storage: storage}
}

func (c *GetLogHandler) Handle(_ context.Context, command *GetLogQuery) (response *GetLogResponse, err error) {
	log, err := c.storage.GetNodeLog(command.NodeId, command.LogPosition)
	if err != nil {
		return nil, err
	}

	response = &GetLogResponse{Log: log}

	return
}
