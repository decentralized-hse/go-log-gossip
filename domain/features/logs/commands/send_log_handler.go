package commands

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/infra/gossip"
	"github.com/decentralized-hse/go-log-gossip/storage"
)

type SendLogHandler struct {
	storage  storage.LogStorage
	gossiper *gossip.Gossiper
}

func NewSendLogHandler(storage storage.LogStorage, gossiper *gossip.Gossiper) *SendLogHandler {
	return &SendLogHandler{storage: storage, gossiper: gossiper}
}

func (c *SendLogHandler) Handle(_ context.Context, command *SendLogCommand) (response *SendLogResponse, err error) {
	log, err := c.storage.GetNodeLog(command.NodeId, command.LogPosition)
	if err != nil {
		return nil, err
	}

	response = &SendLogResponse{Log: log}
	c.gossiper.Request(command.SenderId, gossip.Push, response.Log)
	return
}
