package commands

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/domain/features/dtos"
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

func (c *SendLogHandler) Handle(_ context.Context, command *SendLogCommand) (*SendLogResponse, error) {
	log, _ := c.storage.GetNodeLog(command.NodeId, command.LogPosition)

	response := &SendLogResponse{Log: log}

	if log == nil {
		return response, nil
	}

	dto := dtos.NewLogDTO(log)
	_ = c.gossiper.Request(command.SenderId, gossip.Push, dto)
	return response, nil
}
