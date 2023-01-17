package commands

import (
	"context"
	"github.com/decentralized-hse/go-log-gossip/domain/features/dtos"
	"github.com/decentralized-hse/go-log-gossip/infra/gossip"
	"github.com/decentralized-hse/go-log-gossip/storage"
)

type AddLogHandler struct {
	storage  storage.LogStorage
	gossiper *gossip.Gossiper
}

func NewAddLogHandler(storage storage.LogStorage, gossiper *gossip.Gossiper) *AddLogHandler {
	return &AddLogHandler{storage: storage, gossiper: gossiper}
}

func (c *AddLogHandler) Handle(_ context.Context, command *AddLogCommand) (*AddLogResponse, error) {
	_, requiredLogPosition := c.storage.InsertAt(*command.Log, command.Log.NodeId, command.Log.Position)
	if requiredLogPosition >= 0 {
		requestDto := dtos.NewRequestDTO(command.Log.NodeId, requiredLogPosition)
		_ = c.gossiper.BroadcastMessage(gossip.Pull, requestDto)
	}
	//dto := dtos.NewLogDTO(command.Log)
	//_ = c.gossiper.BroadcastMessage(gossip.Push, dto)

	return &AddLogResponse{}, nil
}
