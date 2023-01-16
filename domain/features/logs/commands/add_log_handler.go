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

func (c *AddLogHandler) Handle(_ context.Context, command *AddLogCommand) (response *AddLogResponse, err error) {
	err, requiredLogPosition := c.storage.InsertAt(command.Log, command.Log.NodeId, command.Log.Position)
	if err != nil && requiredLogPosition >= 0 {
		// TODO: возможно следует обернуть requiredLogPosition в объект
		//_ = c.gossiper.BroadcastMessage(gossip.Pull, requiredLogPosition)
	}
	dto := dtos.NewLogDTO(&command.Log)
	_ = c.gossiper.BroadcastMessage(gossip.Push, dto)
	response = &AddLogResponse{}
	return
}
